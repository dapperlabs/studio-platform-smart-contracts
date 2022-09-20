import PDS from "../../contracts/PDS.cdc"

pub fun main(distId: UInt64): UInt8 {
    return PDS.getDistInfo(distId: distId)!.state.rawValue
}
