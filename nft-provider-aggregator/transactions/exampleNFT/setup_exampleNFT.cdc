import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

/// This transaction sets up the signer's account to hold ExampleNFT NFTs if it hasn't already been configured.
///
transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // Return early if the account already has a collection
        if signer.storage.borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath) != nil {
            return
        }

        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        // Create a new collection and save it to the account storage
        signer.storage.save(<- ExampleNFT.createEmptyCollection(nftType: Type<@ExampleNFT.NFT>()), to: collectionData.storagePath)

        // Create a public capability for the collection
        signer.capabilities.unpublish(ExampleNFT.CollectionPublicPath)
        signer.capabilities.publish(
            signer.capabilities.storage.issue<&ExampleNFT.Collection>(collectionData.storagePath),
            at: collectionData.publicPath
        )
    }
}