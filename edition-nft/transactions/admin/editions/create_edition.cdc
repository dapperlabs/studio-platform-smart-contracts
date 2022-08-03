import EditionNFT from "../../../contracts/EditionNFT.cdc"

transaction(
    metadata: {String: String}
   ) {
    // local variable for the admin reference
    let admin: &EditionNFT.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&EditionNFT.Admin>(from: EditionNFT.AdminStoragePath)
            ?? panic("Could not borrow a reference to the EditionNFT Admin capability")
    }

    execute {
        let id = self.admin.createEdition(
            metadata: metadata
        )

        log("====================================")
        log("New Edition:")
        log("EdiionID: ".concat(id.toString()))
        log("====================================")
    }
}

