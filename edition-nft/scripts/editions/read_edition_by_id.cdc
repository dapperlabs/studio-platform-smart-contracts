import EditionNFT from "EditionNFT"

// This script returns an Edition for an id number, if it exists.

access(all) fun main(editionID: UInt64): EditionNFT.EditionData {
    return EditionNFT.getEditionData(id: editionID)
}