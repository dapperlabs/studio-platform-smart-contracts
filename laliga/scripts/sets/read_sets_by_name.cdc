import Golazos from "../../contracts/Golazos.cdc"

// This script returns a Set struct for the given name,
// if it exists

pub fun main(setName: String): Golazos.SetData {
    return Golazos.getSetDataByName(name: setName)!
}

