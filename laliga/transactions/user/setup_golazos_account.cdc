import NonFungibleToken from "NonFungibleToken"
import Golazos from "Golazos"
import MetadataViews from "MetadataViews"

/// This transaction sets up the signer's account to hold Golazos NFTs if it hasn't already been configured.
///
transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // Return early if the account already has a collection
        if signer.storage.borrow<&Golazos.Collection>(from: Golazos.CollectionStoragePath) != nil {
            return
        }

        // Create a new collection and save it to the account storage
        signer.storage.save(<- Golazos.createEmptyCollection(nftType: Type<@Golazos.NFT>()), to: Golazos.CollectionStoragePath)

        // Create a public capability for the collection
        signer.capabilities.unpublish(Golazos.CollectionPublicPath)
        signer.capabilities.publish(
            signer.capabilities.storage.issue<&Golazos.Collection>(Golazos.CollectionStoragePath),
            at: Golazos.CollectionPublicPath
        )
    }
}