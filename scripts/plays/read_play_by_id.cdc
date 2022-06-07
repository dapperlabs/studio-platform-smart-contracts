import Sport from "../../contracts/Sport.cdc"

// This script returns a Play struct for the given id,
// if it exists

pub fun main(id: UInt64): Sport.PlayData {
    return Sport.getPlayData(id: id)
}

