import Sport from "../../contracts/Sport.cdc"

// This script returns all the Series structs.
// This will eventually be *long*.

pub fun main(): [Sport.SeriesData] {
    let series: [Sport.SeriesData] = []
    var id: UInt64 = 1
    // Note < , as nextSeriesID has not yet been used
    while id < Sport.nextSeriesID {
        series.append(Sport.getSeriesData(id: id))
        id = id + 1
    }
    return series
}

