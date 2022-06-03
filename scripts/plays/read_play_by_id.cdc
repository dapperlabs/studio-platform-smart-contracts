import Golazo from "../../contracts/Golazo.cdc"

// This script returns a Play struct for the given id,
// if it exists

pub fun main(id: UInt64): Golazo.PlayData {
    return Golazo.getPlayData(id: id)
}

