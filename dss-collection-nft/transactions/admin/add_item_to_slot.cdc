import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(slotID: UInt64, id: UInt64) {
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        self.admin.addItemToSlot(
            slotID: slotID,
            id: id
        )
    }
}

