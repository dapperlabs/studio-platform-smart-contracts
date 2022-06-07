import Sport from "../../../contracts/Sport.cdc"

transaction(editionID: UInt64) {
    // local variable for the admin reference
    let admin: &Sport.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&Sport.Admin>(from: Sport.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Sport Admin capability")
    }

    execute {
        let id = self.admin.closeEdition(id: editionID)

        log("====================================")
        log("Closed Edition:")
        log("EditionID: ".concat(id.toString()))
        log("====================================")
    }
}

