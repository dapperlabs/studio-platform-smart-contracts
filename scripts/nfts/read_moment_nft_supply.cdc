import Golazo from "../../contracts/Golazo.cdc"

// This scripts returns the number of Golazo currently in existence.

pub fun main(): UInt64 {    
    return Golazo.totalSupply
}

