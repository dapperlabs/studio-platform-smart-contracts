import Golazos from "../../contracts/Golazos.cdc"

// This script returns all the names for Series.
// These can be related to Series structs via Golazos.getSeriesByName() .

pub fun main(): [String] {
    return Golazos.getAllSeriesNames()
}

