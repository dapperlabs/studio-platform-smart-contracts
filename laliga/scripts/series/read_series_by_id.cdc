import Golazos from "Golazos"

// This script returns a Series struct for the given id,
// if it exists

access(all) fun main(id: UInt64): Golazos.SeriesData {
    return Golazos.getSeriesData(id: id)
}

