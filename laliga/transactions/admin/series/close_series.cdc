import Golazos from "Golazos"

transaction(seriesID: UInt64) {
    // local variable for the admin reference
    let admin: auth(Golazos.Operate) &Golazos.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the Admin resource
        self.admin = signer.storage.borrow<auth(Golazos.Operate) &Golazos.Admin>(from: Golazos.AdminStoragePath)
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

