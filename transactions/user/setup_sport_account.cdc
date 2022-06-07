import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Sport from "../../contracts/Sport.cdc"

// This transaction configures an account to hold Sport NFTs.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&Sport.Collection>(from: Sport.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- Sport.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: Sport.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&Sport.Collection{NonFungibleToken.CollectionPublic, Sport.MomentNFTCollectionPublic}>(
                Sport.CollectionPublicPath,
                target: Sport.CollectionStoragePath
            )
        }
    }
}
