import DSSCollection from 0xf8d6e0586b0a20c7


pub fun main(collectionGroupID: UInt64): DSSCollection.CollectionGroupData {
    return DSSCollection.getCollectionGroupData(id: collectionGroupID)
}