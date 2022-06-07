import Sport from "../../contracts/Sport.cdc"

// This script returns all the Edition structs.
// This will be *long*.

pub fun main(): [Sport.EditionData] {
    let editions: [Sport.EditionData] = []
    var id: UInt64 = 1
    // Note < , as nextEditionID has not yet been used
    while id < Sport.nextEditionID {
        editions.append(Sport.getEditionData(id: id))
        id = id + 1
    }
    return editions
}

