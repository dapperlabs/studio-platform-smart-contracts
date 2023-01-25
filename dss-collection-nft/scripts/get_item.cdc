import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(id: UInt64): DSSCollection.Item {
    return DSSCollection.getItem(id: id)
}