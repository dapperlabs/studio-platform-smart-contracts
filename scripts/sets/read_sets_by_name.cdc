import Golazo from "../../contracts/Golazo.cdc"

// This script returns a Set struct for the given name,
// if it exists

pub fun main(setName: String): Golazo.SetData {
    return Golazo.getSetDataByName(name: setName)
}

