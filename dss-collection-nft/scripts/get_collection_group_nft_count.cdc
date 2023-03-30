import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(collectionGroupId: UInt64): UInt64 {
    let count = DSSCollection.collectionGroupNFTCount[collectionGroupId] ?? 0
    return count
}
