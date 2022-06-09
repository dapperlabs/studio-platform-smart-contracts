import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns a Series struct for the given id,
// if it exists

pub fun main(id: UInt64): DapperSport.SeriesData {
    return DapperSport.getSeriesData(id: id)
}

