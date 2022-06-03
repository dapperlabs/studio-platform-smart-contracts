import Golazo from "../../contracts/Golazo.cdc"

// This script returns all the Series structs.
// This will eventually be *long*.

pub fun main(): [Golazo.SeriesData] {
    let series: [Golazo.SeriesData] = []
    var id: UInt64 = 1
    // Note < , as nextSeriesID has not yet been used
    while id < Golazo.nextSeriesID {
        series.append(Golazo.getSeriesData(id: id))
        id = id + 1
    }
    return series
}

