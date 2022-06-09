import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns a Set struct for the given name,
// if it exists

pub fun main(setName: String): DapperSport.SetData {
    return DapperSport.getSetDataByName(name: setName)
}

