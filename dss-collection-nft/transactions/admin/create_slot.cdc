import DSSCollection from "../../contracts/DSSCollection.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(collectionGroupID: UInt64, logicalOperator: String, required: Bool, metadata: {String: String}) {
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        let typeName = Type<@ExampleNFT.NFT>()
        let id = self.admin.createSlot(
            collectionGroupID: collectionGroupID,
            logicalOperator: logicalOperator,
            required: required,
            typeName: typeName,
            metadata: metadata
        )

        log("====================================")
        log("New Slot:")
        log("ID: ".concat(id.toString()))
        log("====================================")
    }
}