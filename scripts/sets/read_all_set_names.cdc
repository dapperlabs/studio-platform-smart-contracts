import Sport from "../../contracts/Sport.cdc"

// This script returns all the names for Set.
// These can be related to Set structs via Sport.getSetByName() .

pub fun main(): [String] {
    return Sport.getAllSetNames()
}

