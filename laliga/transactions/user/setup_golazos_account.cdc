import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Golazos from "../../contracts/Golazos.cdc"

// This transaction configures an account to hold Golazos NFTs.

transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&Golazos.Collection>(from: Golazos.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- Golazos.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: Golazos.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&Golazos.Collection{NonFungibleToken.CollectionPublic, Golazos.MomentNFTCollectionPublic}>(
                Golazos.CollectionPublicPath,
                target: Golazos.CollectionStoragePath
            )
        }
    }
}
