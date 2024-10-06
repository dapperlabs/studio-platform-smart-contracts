import Escrow from "Escrow"
import AllDay from "AllDay"
import NonFungibleToken from "NonFungibleToken"

// This transaction takes the leaderboardName and nftID and returns it to the correct owner.
transaction(leaderboardName: String, nftID: UInt64, ownerAddress: Address) {
    prepare(signer: auth(BorrowValue) &Account) {
        // Get a reference to the Collection resource in storage.
        let collectionRef = signer.storage.borrow<auth(Escrow.Operate) &Escrow.Collection>(from: Escrow.CollectionStoragePath)
            ?? panic("Could not borrow reference to the Collection resource")

        let depositCap = getAccount(ownerAddress)
            .capabilities.get<&{NonFungibleToken.Collection}>(AllDay.CollectionPublicPath)

        // Call transferNftToCollection function.
        collectionRef.adminTransferNftToCollection(leaderboardName: leaderboardName, nftID: nftID, depositCap: depositCap)
    }

    execute {
        log("Withdrawn NFT from leaderboard")
    }
}
