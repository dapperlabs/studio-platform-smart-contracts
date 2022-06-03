import Golazo from "../../contracts/Golazo.cdc"

// This script returns all the names for Series.
// These can be related to Series structs via Golazo.getSeriesByName() .

pub fun main(): [String] {
    return Golazo.getAllSeriesNames()
}

