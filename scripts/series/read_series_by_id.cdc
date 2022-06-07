import Sport from "../../contracts/Sport.cdc"

// This script returns a Series struct for the given id,
// if it exists

pub fun main(id: UInt64): Sport.SeriesData {
    return Sport.getSeriesData(id: id)
}

