import NonFungibleToken from "NonFungibleToken"
import AllDay from "AllDay"
import Escrow from "Escrow"

transaction(leaderboardName: String, nftID: UInt64) {
    let nft: @{NonFungibleToken.NFT}
    let receiver: Capability<&{NonFungibleToken.CollectionPublic}>
    let collectionPublic: &Escrow.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow a reference to the user's NFT collection as a Provider
        let collectionRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>(
            from: AllDay.CollectionStoragePath
        ) ?? panic("Could not borrow NFT collection reference")

        // Borrow a reference to the user's NFT collection as a Receiver.
        self.receiver = signer.capabilities.get<&{NonFungibleToken.CollectionPublic}>(AllDay.CollectionPublicPath)

        // Withdraw the NFT from the user's collection
        self.nft <- collectionRef.withdraw(withdrawID: nftID)

        // Get the public leaderboard collection
        let escrowAccount = getAccount(0xf8d6e0586b0a20c7)
        self.collectionPublic = escrowAccount
            .capabilities.borrow<&Escrow.Collection>(Escrow.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the public leaderboard collection")
    }

    execute {
        let metadata: {String: String} = {}
        // Add the NFT entry to the leaderboard
        self.collectionPublic.addEntryToLeaderboard(
            nft: <-self.nft,
            leaderboardName: leaderboardName,
            ownerAddress: self.receiver.address,
            metadata: metadata
        )
    }
}
