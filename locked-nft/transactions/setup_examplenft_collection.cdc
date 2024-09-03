import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // Return early if the collection already exists
        if signer.storage.borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath) != nil {
            return
        }

        // Create a new collection and save it to storage
        signer.storage.save(<- ExampleNFT.createEmptyCollection(nftType: Type<@ExampleNFT.NFT>()), to: ExampleNFT.CollectionStoragePath)

        // Create a public capability for the collection
        signer.capabilities.unpublish(ExampleNFT.CollectionPublicPath)
        signer.capabilities.publish(
            signer.capabilities.storage.issue<&ExampleNFT.Collection>(ExampleNFT.CollectionStoragePath),
            at: ExampleNFT.CollectionPublicPath
        )
    }
}