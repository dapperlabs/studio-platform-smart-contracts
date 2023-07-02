import Golazos from "../../contracts/Golazos.cdc"

// This script returns an Edition for an id number, if it exists.

pub fun main(editionID: UInt64): Golazos.EditionData {
    return Golazos.getEditionData(id: editionID)!
}

