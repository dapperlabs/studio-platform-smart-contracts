import NonFungibleToken from "NonFungibleToken"
import Golazos from "Golazos"

// This script returns the size of an account's Golazos collection.

access(all) fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}

