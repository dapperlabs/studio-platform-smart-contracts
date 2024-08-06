import NonFungibleToken from "NonFungibleToken"
import EditionNFT from "EditionNFT"

// This transaction configures an account to hold EditionNFTs.

transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // if the account doesn't already have a collection
        if signer.storage.borrow<&EditionNFT.Collection>(from: EditionNFT.CollectionStoragePath) != nil {
            return
        }

        // Create a new collection and save it to the account storage
        signer.storage.save(<- EditionNFT.createEmptyCollection(nftType: Type<@EditionNFT.NFT>()), to: EditionNFT.CollectionStoragePath)

        // Create a public capability for the collection
        signer.capabilities.unpublish(EditionNFT.CollectionPublicPath)
        signer.capabilities.publish(
            signer.capabilities.storage.issue<&EditionNFT.Collection>(EditionNFT.CollectionStoragePath),
            at: EditionNFT.CollectionPublicPath
        )

    }
}
