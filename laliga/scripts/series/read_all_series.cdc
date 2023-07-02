import Golazos from "../../contracts/Golazos.cdc"

// This script returns all the Series structs.
// This will eventually be *long*.

pub fun main(): [Golazos.SeriesData] {
    let series: [Golazos.SeriesData] = []
    var id: UInt64 = 1
    // Note < , as nextSeriesID has not yet been used
    while id < Golazos.nextSeriesID {
        series.append(Golazos.getSeriesData(id: id)!)
        id = id + 1
    }
    return series
}

