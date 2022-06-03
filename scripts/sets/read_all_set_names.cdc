import Golazo from "../../contracts/Golazo.cdc"

// This script returns all the names for Set.
// These can be related to Set structs via Golazo.getSetByName() .

pub fun main(): [String] {
    return Golazo.getAllSetNames()
}

