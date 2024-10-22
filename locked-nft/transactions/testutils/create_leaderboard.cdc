import Escrow from "Escrow"
import ExampleNFT from "ExampleNFT"
import NonFungibleToken from "NonFungibleToken"

// This transaction takes a name and creates a new leaderboard with that name.
transaction(leaderboardName: String) {
    prepare(signer: auth(BorrowValue) &Account) {
        let collectionRef = signer.storage.borrow<auth(Escrow.Operate) &Escrow.Collection>(from: Escrow.CollectionStoragePath)
            ?? panic("Could not borrow reference to the Collection resource")

        let type = Type<@ExampleNFT.NFT>()

        let newNFTCollection <- ExampleNFT.createEmptyCollection(nftType: Type<@ExampleNFT.NFT>())

        collectionRef.createLeaderboard(name: leaderboardName, nftType: type, collection: <-newNFTCollection)
    }
}
