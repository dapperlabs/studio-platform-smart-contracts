import EditionNFT from "EditionNFT"

transaction(
    metadata: {String: String}
   ) {
    // local variable for the admin reference
    let admin: auth(EditionNFT.Operate) &EditionNFT.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the Admin resource
        self.admin = signer.storage.borrow<auth(EditionNFT.Operate) &EditionNFT.Admin>(from: EditionNFT.AdminStoragePath)
            ?? panic("Could not borrow a reference to the Golazos Admin capability")
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

