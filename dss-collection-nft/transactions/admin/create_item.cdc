import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(itemID: UInt64, points: UInt64, itemType: String) {
    // local variable for the admin reference
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        let id = self.admin.createItem(
            itemID: itemID,
            points: points,
            itemType: itemType
        )

        log("====================================")
        log("New Item:")
        log("ID: ".concat(id.toString()))
        log("====================================")
    }
}