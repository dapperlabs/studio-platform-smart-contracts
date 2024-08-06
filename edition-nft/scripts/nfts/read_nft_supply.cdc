import EditionNFT from "EditionNFT"

// This scripts returns the number of EditionNFT currently in existence.

access(all) fun main(): UInt64 {    
    return EditionNFT.totalSupply
}

