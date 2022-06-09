import DapperSport from "../../contracts/DapperSport.cdc"

// This script returns a Play struct for the given id,
// if it exists

pub fun main(id: UInt64): DapperSport.PlayData {
    return DapperSport.getPlayData(id: id)
}

