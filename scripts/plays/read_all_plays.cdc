import Golazo from "../../contracts/Golazo.cdc"

// This script returns all the Set structs.
// This will eventually be *long*.

pub fun main(): [Golazo.PlayData] {
    let plays: [Golazo.PlayData] = []
    var id: UInt64 = 1
    // Note < , as nextPlayID has not yet been used
    while id < Golazo.nextPlayID {
        plays.append(Golazo.getPlayData(id: id))
        id = id + 1
    }
    return plays
}

