import Golazo from "../../contracts/Golazo.cdc"

// This script returns a Set struct for the given id,
// if it exists

pub fun main(id: UInt64): Golazo.SetData {
    return Golazo.getSetData(id: id)
}

