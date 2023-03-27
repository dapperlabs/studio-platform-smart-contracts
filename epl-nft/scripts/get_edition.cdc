import EnglishPremierLeague from "../../contracts/EnglishPremierLeague.cdc"

pub fun main(editionID: UInt64): EnglishPremierLeague.Edition? {
    return EnglishPremierLeague.getEdition(id: editionID)
}