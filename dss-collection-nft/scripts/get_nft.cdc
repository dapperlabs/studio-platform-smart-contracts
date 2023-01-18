import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"


pub fun main(address: Address, id: UInt64): [AnyStruct] {
    let account = getAccount(address)

    let collectionRef = account.getCapability(DSSCollection.CollectionPublicPath)
        .borrow<&{DSSCollection.DSSCollectionNFTCollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")

    let nft = collectionRef.borrowDSSCollectionNFT(id: id)
        ?? panic("Couldn't borrow DSS Collection NFT")

    return [nft.id, nft.collectionGroupID, nft.serialNumber, nft.completionDate, nft.completedBy]
}