import Golazos from "../../../contracts/Golazos.cdc"

transaction(name: String) {
    // local variable for the admin reference
    let admin: &Golazos.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&Golazos.Admin>(from: Golazos.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Golazos Admin capability")
    }

    execute {
        let id = self.admin.createSeries(
            name: name,
        )

        log("====================================")
        log("New Series: ".concat(name))
        log("SeriesID: ".concat(id.toString()))
        log("====================================")
    }
}

