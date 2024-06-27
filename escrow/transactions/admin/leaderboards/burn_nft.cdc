import Escrow from "Escrow"
import AllDay from "AllDay"

// This transaction takes the leaderboardName and nftID and burns the NFT.
transaction(leaderboardName: String, nftID: UInt64) {
    prepare(signer: auth(BorrowValue) &Account) {
        // Get a reference to the Collection resource in storage.
        let collectionRef = signer.storage.borrow<auth(Escrow.Operate) &Escrow.Collection>(from: Escrow.CollectionStoragePath)
            ?? panic("Could not borrow reference to the Collection resource")

        // Call withdraw function.
        collectionRef.burn(leaderboardName: leaderboardName, nftID: nftID)
    }

    execute {
        log("Burned NFT from leaderboard")
    }
}
