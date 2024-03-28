import PDS from "PDS"

access(all) fun main(distId: UInt64): UInt8 {
    return PDS.getDistInfo(distId: distId)!.state.rawValue
}
