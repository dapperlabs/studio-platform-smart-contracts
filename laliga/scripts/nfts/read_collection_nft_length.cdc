import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Golazos from "../../contracts/Golazos.cdc"

// This script returns the size of an account's Golazos collection.

pub fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = account.getCapability(Golazos.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}

