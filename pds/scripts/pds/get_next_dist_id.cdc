import PDS from "../../contracts/PDS.cdc"

pub fun main(): UInt64 {
    return PDS.nextDistId
}
