import NonFungibleToken from "NonFungibleToken"
import Golazos from "Golazos"

// This script returns the size of an account's Golazos collection.

access(all) fun main(address: Address, id: UInt64): [AnyStruct] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
        ?? panic("Could not borrow capability from public collection")
    
    let nft = collectionRef.borrowMomentNFT(id: id)
        ?? panic("Couldn't borrow momentNFT")

    return [nft.id, nft.editionID, nft.serialNumber, nft.mintingDate]
}

