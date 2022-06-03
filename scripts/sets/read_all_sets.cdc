import Golazo from "../../contracts/Golazo.cdc"

// This script returns all the Set structs.
// This will eventually be *long*.

pub fun main(): [Golazo.SetData] {
    let sets: [Golazo.SetData] = []
    var id: UInt64 = 1
    // Note < , as nextSetID has not yet been used
    while id < Golazo.nextSetID {
        sets.append(Golazo.getSetData(id: id))
        id = id + 1
    }
    return sets
}

