import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns a Series struct for the given name,
// if it exists

pub fun main(seriesName: String): DapperSport.SeriesData {
    return DapperSport.getSeriesDataByName(name: seriesName)
}

