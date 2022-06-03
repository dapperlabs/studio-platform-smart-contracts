import Golazo from "../../contracts/Golazo.cdc"

// This script returns an Edition for an id number, if it exists.

pub fun main(editionID: UInt64): Golazo.EditionData {
    return Golazo.getEditionData(id: editionID)
}

