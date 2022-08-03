import EditionNFT from "../../../contracts/EditionNFT.cdc"

transaction(editionID: UInt64) {
    // local variable for the admin reference
    let admin: &EditionNFT.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&EditionNFT.Admin>(from: EditionNFT.AdminStoragePath)
            ?? panic("Could not borrow a reference to the EditionNFT Admin capability")
    }

    execute {
        let id = self.admin.closeEdition(id: editionID)

        log("====================================")
        log("Closed Edition:")
        log("EditionID: ".concat(id.toString()))
        log("====================================")
    }
}

