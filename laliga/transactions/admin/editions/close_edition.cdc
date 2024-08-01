import Golazos from "Golazos"

transaction(editionID: UInt64) {
    // local variable for the admin reference
    let admin: &Golazos.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the Admin resource
        self.admin = signer.storage.borrow<&Golazos.Admin>(from: Golazos.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Golazos Admin capability")
    }

    execute {
        let id = self.admin.closeEdition(id: editionID)

        log("====================================")
        log("Closed Edition:")
        log("EditionID: ".concat(id.toString()))
        log("====================================")
    }
}

