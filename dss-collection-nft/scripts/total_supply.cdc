import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(): UInt64 {
    return DSSCollection.totalSupply
}