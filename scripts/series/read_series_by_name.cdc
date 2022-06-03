import Golazo from "../../contracts/Golazo.cdc"

// This script returns a Series struct for the given name,
// if it exists

pub fun main(seriesName: String): Golazo.SeriesData {
    return Golazo.getSeriesDataByName(name: seriesName)
}

