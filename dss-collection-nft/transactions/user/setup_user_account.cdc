import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&DSSCollection.Collection>(from: DSSCollection.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- DSSCollection.createEmptyCollection()
            
            // save it to the account
            signer.save(<-collection, to: DSSCollection.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&DSSCollection.Collection{NonFungibleToken.CollectionPublic, DSSCollection.DSSCollectionNFTCollectionPublic}>(
                DSSCollection.CollectionPublicPath,
                target: DSSCollection.CollectionStoragePath
            )
        }
    }
}
