import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(
    itemID: String,
    points: UInt64,
    itemType: String,
    comparator: String,
    slotID: UInt64
) {
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        self.admin.createItemInSlot(
            itemID: itemID,
            points: points,
            itemType: itemType,
            comparator: comparator,
            slotID: slotID
        )
    }
}

