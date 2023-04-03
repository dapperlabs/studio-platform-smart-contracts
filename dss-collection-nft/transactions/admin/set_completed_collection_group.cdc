import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(collectionGroupID: UInt64, userAddress: Address, nftIDs: [UInt64]) {
    // Local variable for the admin reference
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        // Borrow a reference to the Admin resource
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        // Call the completedCollectionGroup function using the borrowed admin reference
        self.admin.completedCollectionGroup(
            collectionGroupID: collectionGroupID,
            userAddress: userAddress,
            nftIDs: nftIDs
        )

        log("====================================")
        log("Updated Completed Collection Group:")
        log("CollectionID: ".concat(collectionGroupID.toString()))
        log("UserAddress: ".concat(userAddress.toString()))
        log("NFT IDs: [")
        for nftID in nftIDs {
            log("  ".concat(nftID.toString()))
        }
        log("]")
        log("====================================")
    }
}
