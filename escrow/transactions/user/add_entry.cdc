import NonFungibleToken from "NonFungibleToken"
import AllDay from "AllDay"
import Escrow from "Escrow"

transaction(leaderboardName: String, nftID: UInt64) {
    let nft: @{NonFungibleToken.NFT}
    let receiverCollection: Capability<&{NonFungibleToken.Collection}>
    let escrowCollection: &Escrow.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow a reference to the user's NFT collection as a Provider
        let collectionRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
            from: AllDay.CollectionStoragePath
        ) ?? panic("Could not borrow NFT collection reference")

        // Borrow a reference to the user's NFT collection as a Receiver.
        self.receiverCollection = signer.capabilities.get<&{NonFungibleToken.Collection}>(AllDay.CollectionPublicPath)

        // Withdraw the NFT from the user's collection
        self.nft <- collectionRef.withdraw(withdrawID: nftID)

        // Get the public leaderboard collection
        let escrowAccount = getAccount(0xf8d6e0586b0a20c7)
        self.escrowCollection = escrowAccount
            .capabilities.borrow<&Escrow.Collection>(Escrow.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the public leaderboard collection")
    }

    execute {
        // Add the NFT entry to the leaderboard
        self.escrowCollection.addEntryToLeaderboard(
            nft: <-self.nft,
            leaderboardName: leaderboardName,
            ownerAddress: self.receiverCollection.address,
        )
    }
}
