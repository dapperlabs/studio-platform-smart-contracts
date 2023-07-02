import Golazos from "../../contracts/Golazos.cdc"

// This script returns a Set struct for the given id,
// if it exists

pub fun main(id: UInt64): Golazos.SetData {
    return Golazos.getSetData(id: id)!
}

