import NonFungibleToken from "NonFungibleToken"
import Golazos from "Golazos"

// This script returns an array of all the NFT IDs in an account's collection.

access(all) fun main(address: Address): [UInt64] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs()
}

