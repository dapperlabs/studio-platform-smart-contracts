import Golazos from "../../../contracts/Golazos.cdc"

transaction(seriesID: UInt64) {
    // local variable for the admin reference
    let admin: &Golazos.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&Golazos.Admin>(from: Golazos.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Golazos Admin capability")
    }

    execute {
        let id = self.admin.closeSeries(id: seriesID)

        log("====================================")
        log("Closed Series:")
        log("SeriesID: ".concat(id.toString()))
        log("====================================")
    }
}

