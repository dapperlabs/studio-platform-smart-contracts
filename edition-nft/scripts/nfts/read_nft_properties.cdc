import NonFungibleToken from "NonFungibleToken"
import EditionNFT from "EditionNFT"

// This script returns the size of an account's EditionNFT collection.

access(all) fun main(address: Address, id: UInt64): [AnyStruct] {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&EditionNFT.Collection>(EditionNFT.CollectionPublicPath)
        ?? panic("Could not borrow capability from public collection")
    
    let nft = collectionRef.borrowEditionNFT(id: id)
        ?? panic("Couldn't borrow momentNFT")

    return [nft.id, nft.editionID]
}

