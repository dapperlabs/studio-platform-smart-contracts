import Golazo from "../../contracts/Golazo.cdc"

// This script returns all the Edition structs.
// This will be *long*.

pub fun main(): [Golazo.EditionData] {
    let editions: [Golazo.EditionData] = []
    var id: UInt64 = 1
    // Note < , as nextEditionID has not yet been used
    while id < Golazo.nextEditionID {
        editions.append(Golazo.getEditionData(id: id))
        id = id + 1
    }
    return editions
}

