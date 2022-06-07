import Sport from "../../contracts/Sport.cdc"

// This script returns a Set struct for the given id,
// if it exists

pub fun main(id: UInt64): Sport.SetData {
    return Sport.getSetData(id: id)
}

