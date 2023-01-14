import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(nftID: UInt64, collectionGroupID: UInt64) {
    // local variable for the admin reference
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        self.admin.addNFTToCollectionGroup(
            nftID: nftID, 
            collectionGroupID: collectionGroupID
        )

        log("====================================")
        log("Added NFT to Collection Group: ".concat(nftID.toString()))
        log("CollectionGroupID: ".concat(collectionGroupID.toString()))
        log("====================================")
    }
}