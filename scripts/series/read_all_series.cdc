import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns all the Series structs.
// This will eventually be *long*.

pub fun main(): [DapperSport.SeriesData] {
    let series: [DapperSport.SeriesData] = []
    var id: UInt64 = 1
    // Note < , as nextSeriesID has not yet been used
    while id < DapperSport.nextSeriesID {
        series.append(DapperSport.getSeriesData(id: id))
        id = id + 1
    }
    return series
}

