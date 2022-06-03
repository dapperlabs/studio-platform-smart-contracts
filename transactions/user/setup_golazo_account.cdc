import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Golazo from "../../contracts/Golazo.cdc"

// This transaction configures an account to hold Golazo NFTs.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&Golazo.Collection>(from: Golazo.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- Golazo.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: Golazo.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&Golazo.Collection{NonFungibleToken.CollectionPublic, Golazo.MomentNFTCollectionPublic}>(
                Golazo.CollectionPublicPath,
                target: Golazo.CollectionStoragePath
            )
        }
    }
}
