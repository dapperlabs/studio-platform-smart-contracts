import Golazos from "../../contracts/Golazos.cdc"

// This script returns a Series struct for the given id,
// if it exists

pub fun main(id: UInt64): Golazos.SeriesData {
    return Golazos.getSeriesData(id: id)!
}

