import Sport from "../../contracts/Sport.cdc"

// This script returns a Series struct for the given name,
// if it exists

pub fun main(seriesName: String): Sport.SeriesData {
    return Sport.getSeriesDataByName(name: seriesName)
}

