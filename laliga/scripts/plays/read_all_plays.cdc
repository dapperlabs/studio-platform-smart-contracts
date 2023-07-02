import Golazos from "../../contracts/Golazos.cdc"

// This script returns all the Set structs.
// This will eventually be *long*.

pub fun main(): [Golazos.PlayData] {
    let plays: [Golazos.PlayData] = []
    var id: UInt64 = 1
    // Note < , as nextPlayID has not yet been used
    while id < Golazos.nextPlayID {
        plays.append(Golazos.getPlayData(id: id)!)
        id = id + 1
    }
    return plays
}

