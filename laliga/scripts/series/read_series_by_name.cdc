import Golazos from "../../contracts/Golazos.cdc"

// This script returns a Series struct for the given name,
// if it exists

pub fun main(seriesName: String): Golazos.SeriesData {
    return Golazos.getSeriesDataByName(name: seriesName)!
}

