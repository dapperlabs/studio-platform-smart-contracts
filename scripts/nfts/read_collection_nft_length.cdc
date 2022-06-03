import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Golazo from "../../contracts/Golazo.cdc"

// This script returns the size of an account's Golazo collection.

pub fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = account.getCapability(Golazo.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}

