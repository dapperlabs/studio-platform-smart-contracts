import DSSCollection from "../../contracts/DSSCollection.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(collectionGroupID: UInt64, logicalOperator: String) {
    // local variable for the admin reference
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        let typeName = Type<@ExampleNFT.NFT>()
        let id = self.admin.createSlot(
            collectionGroupID: collectionGroupID,
            logicalOperator: logicalOperator,
            typeName: typeName
        )

        log("====================================")
        log("New Slot:")
        log("ID: ".concat(id.toString()))
        log("====================================")
    }
}