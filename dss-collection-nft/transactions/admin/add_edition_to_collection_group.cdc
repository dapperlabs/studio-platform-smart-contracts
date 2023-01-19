import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(collectionGroupID: UInt64, editionID: UInt64) {
    // local variable for the admin reference
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        self.admin.addEditionToCollectionGroup(
            collectionGroupID: collectionGroupID,
            editionID: editionID
        )

        log("====================================")
        log("CollectionGroupID: ".concat(collectionGroupID.toString()))
        log("Added Edition to Collection Group: ".concat(editionID.toString()))
        log("====================================")
    }
}

