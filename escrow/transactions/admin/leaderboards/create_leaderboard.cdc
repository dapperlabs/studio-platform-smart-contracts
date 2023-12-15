import Escrow from "../../../contracts/AllDay.cdc"
import AllDay from "../../../contracts/AllDay.cdc"
import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"

// This transaction takes a name and creates a new leaderboard with that name.
transaction(leaderboardName: String) {
    prepare(signer: AuthAccount) {
        let collectionRef = signer.borrow<&Escrow.Collection>(from: Escrow.CollectionStoragePath)
            ?? panic("Could not borrow reference to the Collection resource")

        let type = Type<@AllDay.NFT>()

        let newNFTCollection <- AllDay.createEmptyCollection() as! @NonFungibleToken.Collection

        collectionRef.createLeaderboard(name: leaderboardName, nftType: type, collection: <-newNFTCollection)
    }
}
