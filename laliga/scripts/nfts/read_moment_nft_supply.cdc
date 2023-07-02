import Golazos from "../../contracts/Golazos.cdc"

// This scripts returns the number of Golazos currently in existence.

pub fun main(): UInt64 {    
    return Golazos.totalSupply
}

