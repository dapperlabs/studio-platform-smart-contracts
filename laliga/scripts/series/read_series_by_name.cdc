import Golazos from "Golazos"

// This script returns a Series struct for the given name,
// if it exists

access(all) fun main(seriesName: String): Golazos.SeriesData {
    return Golazos.getSeriesDataByName(name: seriesName)!
}

