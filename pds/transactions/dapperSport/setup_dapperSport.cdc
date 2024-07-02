import NonFungibleToken from "NonFungibleToken"
import DapperSport from "DapperSport"

// This transaction configures an account to hold DapperSport NFTs.

transaction {
    prepare(signer: auth(Storage, Capabilities) &Account) {
        // if the account doesn't already have a collection
        if signer.storage.borrow<&DapperSport.Collection>(from: DapperSport.CollectionStoragePath) == nil {

            // create a new empty collection
            let collection <- DapperSport.createEmptyCollection(nftType: Type<@DapperSport.NFT>())

            // save it to the account
            signer.storage.save(<-collection, to: DapperSport.CollectionStoragePath)

            // create a public capability for the collection
            signer.capabilities.publish(
                signer.capabilities.storage.issue<&DapperSport.Collection>(DapperSport.CollectionStoragePath),
                at: DapperSport.CollectionPublicPath
            )
        }
    }
}