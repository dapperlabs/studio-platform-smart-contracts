import Golazos from "Golazos"

// This script returns all the names for Series.
// These can be related to Series structs via Golazos.getSeriesByName() .

access(all) fun main(): [String] {
    return Golazos.getAllSeriesNames()
}

