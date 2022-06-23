import DapperSport from "../../contracts/DapperSport.cdc"

// This scripts returns the number of DapperSport currently in existence.

pub fun main(): UInt64 {    
    return DapperSport.totalSupply
}

