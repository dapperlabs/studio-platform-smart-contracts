import Sport from "../../contracts/Sport.cdc"

// This script returns all the Set structs.
// This will eventually be *long*.

pub fun main(): [Sport.SetData] {
    let sets: [Sport.SetData] = []
    var id: UInt64 = 1
    // Note < , as nextSetID has not yet been used
    while id < Sport.nextSetID {
        sets.append(Sport.getSetData(id: id))
        id = id + 1
    }
    return sets
}

