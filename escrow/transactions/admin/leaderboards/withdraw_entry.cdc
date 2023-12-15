import Escrow from "../../../contracts/AllDay.cdc"
import AllDay from "../../../contracts/AllDay.cdc"
import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"

// This transaction takes the leaderboardName and nftID and returns it to the correct owner.
transaction(leaderboardName: String, nftID: UInt64, ownerAddress: Address) {
    prepare(signer: AuthAccount) {
        // Get a reference to the Collection resource in storage.
        let collectionRef = signer.borrow<&Escrow.Collection>(from: Escrow.CollectionStoragePath)
            ?? panic("Could not borrow reference to the Collection resource")

        let depositCap = getAccount(ownerAddress)
            .getCapability<&{NonFungibleToken.CollectionPublic}>(AllDay.CollectionPublicPath)

        // Call transferNftToCollection function.
        collectionRef.transferNftToCollection(leaderboardName: leaderboardName, nftID: nftID, depositCap: depositCap)
    }

    execute {
        log("Withdrawn NFT from leaderboard")
    }
}
