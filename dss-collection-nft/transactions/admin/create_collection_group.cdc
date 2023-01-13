import DSSCollection from "../../contracts/DSSCollection.cdc"


transaction(name: String, productPublicPath: PublicPath) {
    // local variable for the admin reference
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        // borrow a reference to the Admin resource
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        let id = self.admin.createCollectionGroup(
            name: name,
            productPublicPath: productPublicPath,
            startTime: nil,
            endTime: nil,
            timeBound: false
        )

        log("====================================")
        log("New Collection Group:")
        log("CollectionGroupID: ".concat(id.toString()))
        log("====================================")
    }
}