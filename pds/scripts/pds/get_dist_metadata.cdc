import PDS from "../../contracts/PDS.cdc"

pub fun main(distId: UInt64): {String: String} {
    return PDS.getDistInfo(distId: distId)!.metadata
}
