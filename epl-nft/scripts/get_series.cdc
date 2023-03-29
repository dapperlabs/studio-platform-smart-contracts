import EnglishPremierLeague from "./EnglishPremierLeague.cdc"

pub fun main(seriesID: UInt64): EnglishPremierLeague.Series {
    return EnglishPremierLeague.getSeries(id: seriesID)!
}