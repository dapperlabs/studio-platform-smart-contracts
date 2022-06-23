import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns all the names for Series.
// These can be related to Series structs via DapperSport.getSeriesByName() .

pub fun main(): [String] {
    return DapperSport.getAllSeriesNames()
}

