import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import EditionNFT from "../../contracts/EditionNFT.cdc"

// This script returns the size of an account's AllDay collection.

pub fun main(address: Address, id: UInt64): [AnyStruct] {
    let account = getAccount(address)

    let collectionRef = account.getCapability(EditionNFT.CollectionPublicPath)
        .borrow<&{EditionNFT.EditionNFTCollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")
    
    let nft = collectionRef.borrowEditionNFT(id: id)
        ?? panic("Couldn't borrow momentNFT")

    return [nft.id, nft.editionID]
}

