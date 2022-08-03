import EditionNFT from "../../contracts/EditionNFT.cdc"

// This script returns an Edition for an id number, if it exists.

pub fun main(editionID: UInt64): EditionNFT.EditionData {
    return EditionNFT.getEditionData(id: editionID)
}