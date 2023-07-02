import Golazos from "../../contracts/Golazos.cdc"

// This script returns all the names for Set.
// These can be related to Set structs via Golazos.getSetByName() .

pub fun main(): [String] {
    return Golazos.getAllSetNames()
}

