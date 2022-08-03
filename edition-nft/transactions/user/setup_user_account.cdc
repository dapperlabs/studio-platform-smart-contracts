import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import EditionNFT from "../../contracts/EditionNFT.cdc"

// This transaction configures an account to hold EditionNFT NFTs.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&EditionNFT.Collection>(from: EditionNFT.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- EditionNFT.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: EditionNFT.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&EditionNFT.Collection{NonFungibleToken.CollectionPublic, EditionNFT.EditionNFTCollectionPublic}>(
                EditionNFT.CollectionPublicPath,
                target: EditionNFT.CollectionStoragePath
            )
        }
    }
}
