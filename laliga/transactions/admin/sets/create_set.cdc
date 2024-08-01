import Golazos from "Golazos"

transaction(name: String) {
    // local variable for the admin reference
    let admin: &Golazos.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the Admin resource
        self.admin = signer.storage.borrow<&Golazos.Admin>(from: Golazos.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Golazos Admin capability")
    }

    execute {
        let id = self.admin.createSet(
            name: name,
        )

        log("====================================")
        log("New Set: ".concat(name))
        log("SetID: ".concat(id.toString()))
        log("====================================")
    }
}

