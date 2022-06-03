import Golazo from "../../contracts/Golazo.cdc"

// This script returns a Series struct for the given id,
// if it exists

pub fun main(id: UInt64): Golazo.SeriesData {
    return Golazo.getSeriesData(id: id)
}

