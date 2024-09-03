import Escrow from "Escrow"
import AllDay from "AllDay"
import NonFungibleToken from "NonFungibleToken"

// This transaction takes a name and creates a new leaderboard with that name.
transaction(leaderboardName: String) {
    prepare(signer: auth(BorrowValue) &Account) {
        let collectionRef = signer.storage.borrow<auth(Escrow.Operate) &Escrow.Collection>(from: Escrow.CollectionStoragePath)
            ?? panic("Could not borrow reference to the Collection resource")

        let type = Type<@AllDay.NFT>()

        let newNFTCollection <- AllDay.createEmptyCollection(nftType: Type<@AllDay.NFT>())

        collectionRef.createLeaderboard(name: leaderboardName, nftType: type, collection: <-newNFTCollection)
    }
}
