import Sport from "../../contracts/Sport.cdc"

// This script returns all the names for Series.
// These can be related to Series structs via Sport.getSeriesByName() .

pub fun main(): [String] {
    return Sport.getAllSeriesNames()
}

