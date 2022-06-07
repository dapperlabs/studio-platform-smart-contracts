import Sport from "../../contracts/Sport.cdc"

// This script returns a Set struct for the given name,
// if it exists

pub fun main(setName: String): Sport.SetData {
    return Sport.getSetDataByName(name: setName)
}

