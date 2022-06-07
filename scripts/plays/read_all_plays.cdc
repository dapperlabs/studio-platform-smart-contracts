import Sport from "../../contracts/Sport.cdc"

// This script returns all the Set structs.
// This will eventually be *long*.

pub fun main(): [Sport.PlayData] {
    let plays: [Sport.PlayData] = []
    var id: UInt64 = 1
    // Note < , as nextPlayID has not yet been used
    while id < Sport.nextPlayID {
        plays.append(Sport.getPlayData(id: id))
        id = id + 1
    }
    return plays
}

