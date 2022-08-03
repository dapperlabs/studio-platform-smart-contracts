import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns all the names for Set.
// These can be related to Set structs via DapperSport.getSetByName() .

pub fun main(): [String] {
    return DapperSport.getAllSetNames()
}

