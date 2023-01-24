import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(itemID: UInt64): DSSCollection.ItemData {
    return DSSCollection.getItemData(id: itemID)
}