import Sport from "../../contracts/Sport.cdc"

// This scripts returns the number of Sport currently in existence.

pub fun main(): UInt64 {    
    return Sport.totalSupply
}

