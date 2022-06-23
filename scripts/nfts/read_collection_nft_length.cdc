import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns the size of an account's DapperSport collection.

pub fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = account.getCapability(DapperSport.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    return collectionRef.getIDs().length
}

