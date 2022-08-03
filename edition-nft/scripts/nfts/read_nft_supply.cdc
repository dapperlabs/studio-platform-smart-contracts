import EditionNFT from "../../contracts/EditionNFT.cdc"

// This scripts returns the number of EditionNFT currently in existence.

pub fun main(): UInt64 {    
    return EditionNFT.totalSupply
}

