import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(slotID: UInt64): DSSCollection.SlotData {
    return DSSCollection.getSlotData(id: slotID)
}