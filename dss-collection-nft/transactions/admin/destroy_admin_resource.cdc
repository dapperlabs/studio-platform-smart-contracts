import DSSCollection from "../../contracts/DSSCollection.cdc"


transaction() {
    let adminResource: @DSSCollection.Admin

    prepare(signer: AuthAccount) {
        self.adminResource <- signer.load<@DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not load admin resource")
    }

    execute {
        destroy self.adminResource
    }
}