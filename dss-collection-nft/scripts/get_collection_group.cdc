import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(collectionGroupID: UInt64): DSSCollection.CollectionGroupData {
    return DSSCollection.getCollectionGroupData(id: collectionGroupID)
}