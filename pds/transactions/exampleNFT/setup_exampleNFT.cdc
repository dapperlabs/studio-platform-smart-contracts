import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

// This transaction configures an account to hold ExampleNFT NFTs.

transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // if the account doesn't already have a collection
        if signer.storage.borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- ExampleNFT.createEmptyCollection(nftType: Type<@ExampleNFT.NFT>())

            // save it to the account
            signer.storage.save(<-collection, to: ExampleNFT.CollectionStoragePath)

            // create a public capability for the collection
            signer.capabilities.publish(
                signer.capabilities.storage.issue<&ExampleNFT.Collection>(ExampleNFT.CollectionStoragePath),
                at: ExampleNFT.CollectionPublicPath
            )
        }
    }
}