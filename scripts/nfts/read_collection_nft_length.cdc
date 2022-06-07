import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Sport from "../../contracts/Sport.cdc"

// This script returns the size of an account's Sport collection.

pub fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = account.getCapability(Sport.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}

