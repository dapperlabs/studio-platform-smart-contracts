import Golazos from "../../contracts/Golazos.cdc"

// This script returns a Play struct for the given id,
// if it exists

pub fun main(id: UInt64): Golazos.PlayData {
    return Golazos.getPlayData(id: id)!
}

