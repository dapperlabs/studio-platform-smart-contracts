import Sport from "../../contracts/Sport.cdc"

// This script returns an Edition for an id number, if it exists.

pub fun main(editionID: UInt64): Sport.EditionData {
    return Sport.getEditionData(id: editionID)
}

