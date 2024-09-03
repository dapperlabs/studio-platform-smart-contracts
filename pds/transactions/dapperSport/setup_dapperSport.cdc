import NonFungibleToken from 0x{{.NonFungibleToken}}
import DapperSport from 0x{{.DapperSport}}

// This transaction configures an account to hold DapperSport NFTs.
transaction {
    prepare(signer: AuthAccount) {
        // if the account doesn't already have a collection
        if signer.borrow<&DapperSport.Collection>(from: DapperSport.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- DapperSport.createEmptyCollection()

            // save it to the account
            signer.save(<-collection, to: DapperSport.CollectionStoragePath)

            // create a public capability for the collection
            signer.link<&DapperSport.Collection{NonFungibleToken.CollectionPublic, DapperSport.MomentNFTCollectionPublic}>(
                DapperSport.CollectionPublicPath,
                target: DapperSport.CollectionStoragePath
            )
        }
    }
}