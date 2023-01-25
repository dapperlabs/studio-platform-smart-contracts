import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction(name: String, description: String, typeName: String, endTime: UFix64?) {
    let admin: &DSSCollection.Admin

    prepare(signer: AuthAccount) {
        self.admin = signer.borrow<&DSSCollection.Admin>(from: DSSCollection.AdminStoragePath)
            ?? panic("Could not borrow a reference to the DSSCollection Admin capability")
    }

    execute {
        let id = self.admin.createCollectionGroup(
            name: name,
            description: description,
            typeName: typeName,
            endTime: endTime
        )

        log("====================================")
        log("New Collection Group:")
        log("CollectionGroupID: ".concat(id.toString()))
        log("====================================")
    }
}