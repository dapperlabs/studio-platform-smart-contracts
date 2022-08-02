import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns all the Edition structs.
// This will be *long*.

pub fun main(): [DapperSport.EditionData] {
    let editions: [DapperSport.EditionData] = []
    var id: UInt64 = 1
    // Note < , as nextEditionID has not yet been used
    while id < DapperSport.nextEditionID {
        editions.append(DapperSport.getEditionData(id: id)!)
        id = id + 1
    }
    return editions
}

