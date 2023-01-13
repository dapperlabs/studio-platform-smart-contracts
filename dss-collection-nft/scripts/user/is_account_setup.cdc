import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

// Check to see if an account looks like it has been set up to hold DSSCollections.

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            DSSCollection.DSSCollectionNFTCollectionPublic
        }>(DSSCollection.CollectionPublicPath)
        != nil
}

