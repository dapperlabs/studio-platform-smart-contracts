import EnglishPremierLeague from "../../contracts/EnglishPremierLeague.cdc"

pub fun main(setID: UInt64): EnglishPremierLeague.Set? {
    return EnglishPremierLeague.getSet(id: setID)
}