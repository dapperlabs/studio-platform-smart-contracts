import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import NFTLocker from "NFTLocker"
import Escrow from "Escrow"

/// This transaction unlocks the NFT with provided ID and adds to an escrow leaderboard, unlocking them if necessary.
///
transaction(leaderboardName: String, nftID: UInt64) {
    let ownerAddress: Address
    let collectionRef: auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection
    let collectionPublic: &Escrow.Collection
    let userLockerCollection: auth(NFTLocker.Operate) &NFTLocker.Collection?

    prepare(owner: auth(Storage, Capabilities) &Account) {
        // Borrow a reference to the user's NFT collection as a Provider
        self.collectionRef = owner.storage
            .borrow<auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // Save the owner's address
        self.ownerAddress = owner.address

        // Extract escrow address from contract import
        let escrowAdress = Address.fromString("0x".concat(Type<Escrow>().identifier.slice(from: 2, upTo: 18)))
            ?? panic("Could not convert the address")

        // let escrowAccount = getAccount({{0xEscrowAddress}})
        self.collectionPublic = getAccount(escrowAdress).capabilities.borrow<&Escrow.Collection>(Escrow.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the public leaderboard collection")

        // Borrow a reference to the user's NFTLocker collection
        self.userLockerCollection = owner.storage
            .borrow<auth(NFTLocker.Operate) &NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
    }

    execute {
        // Prepare the NFT type
        let nftType: Type = Type<@ExampleNFT.NFT>()

        // Add NFT to the leaderboard, unlocking it if necessary
        if self.userLockerCollection != nil && NFTLocker.getNFTLockerDetails(id: nftID, nftType: nftType) != nil {
            // Unlock the NFT normally if it has met the unlock conditions, otherwise force unlock (depositing to escrow allows bypassing the unlock conditions)
            if NFTLocker.canUnlockToken(id: nftID, nftType: nftType) {
                self.collectionPublic.addEntryToLeaderboard(
                    nft: <- self.userLockerCollection!.unlock(id: nftID, nftType: nftType),
                    leaderboardName: leaderboardName,
                    ownerAddress: self.ownerAddress,
                )
            } else {
                self.userLockerCollection!.unlockWithAuthorizedDeposit(
                    id: nftID,
                    nftType: nftType,
                    receiverName: "add-entry-to-escrow-leaderboard",
                    passThruParams: {"leaderboardName": leaderboardName},
                )
            }
        } else {
            self.collectionPublic.addEntryToLeaderboard(
                nft: <- self.collectionRef.withdraw(withdrawID: nftID),
                leaderboardName: leaderboardName,
                ownerAddress: self.ownerAddress,
            )
        }
    }
}
