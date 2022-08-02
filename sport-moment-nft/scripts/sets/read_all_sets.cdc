import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns all the Set structs.
// This will eventually be *long*.

pub fun main(): [DapperSport.SetData] {
    let sets: [DapperSport.SetData] = []
    var id: UInt64 = 1
    // Note < , as nextSetID has not yet been used
    while id < DapperSport.nextSetID {
        sets.append(DapperSport.getSetData(id: id)!)
        id = id + 1
    }
    return sets
}

