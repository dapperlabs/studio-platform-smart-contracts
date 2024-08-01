import Golazos from "Golazos"

// This script returns all the Set structs.
// This will eventually be *long*.

access(all) fun main(): [Golazos.SetData] {
    let sets: [Golazos.SetData] = []
    var id: UInt64 = 1
    // Note < , as nextSetID has not yet been used
    while id < Golazos.nextSetID {
        sets.append(Golazos.getSetData(id: id)!)
        id = id + 1
    }
    return sets
}

