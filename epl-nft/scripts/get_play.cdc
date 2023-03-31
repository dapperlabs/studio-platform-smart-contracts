import EnglishPremierLeague from "./EnglishPremierLeague.cdc"

pub fun main(playID: UInt64): EnglishPremierLeague.Play {
    return EnglishPremierLeague.getPlay(id: playID)!
}