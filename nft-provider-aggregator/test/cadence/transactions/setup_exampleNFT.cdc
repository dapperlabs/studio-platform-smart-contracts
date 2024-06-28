import NonFungibleToken from "NonFungibleToken"
import AltExampleNFT from "AltExampleNFT"
import MetadataViews from "MetadataViews"

/// This transaction sets up the signer's account to hold AltExampleNFT NFTs if it hasn't already been configured.
///
transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // Return early if the account already has a collection
        if signer.storage.borrow<&AltExampleNFT.Collection>(from: AltExampleNFT.CollectionStoragePath) != nil {
            return
        }

        let collectionData = AltExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        // Create a new collection and save it to the account storage
        signer.storage.save(<- AltExampleNFT.createEmptyCollection(nftType: Type<@AltExampleNFT.NFT>()), to: collectionData.storagePath)

        // Create a public capability for the collection
        signer.capabilities.unpublish(AltExampleNFT.CollectionPublicPath)
        signer.capabilities.publish(
            signer.capabilities.storage.issue<&AltExampleNFT.Collection>(collectionData.storagePath),
            at: collectionData.publicPath
        )
    }
}