import Golazos from "Golazos"

// This script returns all the Edition structs.
// This will be *long*.

access(all) fun main(): [Golazos.EditionData] {
    let editions: [Golazos.EditionData] = []
    var id: UInt64 = 1
    // Note < , as nextEditionID has not yet been used
    while id < Golazos.nextEditionID {
        editions.append(Golazos.getEditionData(id: id)!)
        id = id + 1
    }
    return editions
}

