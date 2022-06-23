import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns the size of an account's DapperSport collection.

pub fun main(address: Address, id: UInt64): [AnyStruct] {
    let account = getAccount(address)

    let collectionRef = account.getCapability(DapperSport.CollectionPublicPath)
        .borrow<&{DapperSport.MomentNFTCollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    let nft = collectionRef.borrowMomentNFT(id: id)
        ?? panic("Couldn't borrow momentNFT")

    return [nft.id, nft.editionID, nft.serialNumber, nft.mintingDate]
}

