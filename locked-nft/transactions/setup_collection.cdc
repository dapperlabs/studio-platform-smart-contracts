import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"

transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // Return early if the collection already exists
        if signer.storage.borrow<&NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath) != nil {
            return
        }

        // Create a new collection and save it to storage
        signer.storage.save(<- NFTLocker.createEmptyCollection(), to: NFTLocker.CollectionStoragePath)

        // Create a public capability for the collection
        signer.capabilities.unpublish(NFTLocker.CollectionPublicPath)
        signer.capabilities.publish(
            signer.capabilities.storage.issue<&NFTLocker.Collection>(NFTLocker.CollectionStoragePath),
            at: NFTLocker.CollectionPublicPath
        )
    }
}