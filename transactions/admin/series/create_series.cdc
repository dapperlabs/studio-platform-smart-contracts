import Sport from "../../../contracts/Sport.cdc"

transaction(name: String) {
    // local variable for the admin reference
    let admin: &Sport.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&Sport.Admin>(from: Sport.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Sport Admin capability")
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

