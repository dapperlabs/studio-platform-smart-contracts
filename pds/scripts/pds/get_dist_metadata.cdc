import PDS from "PDS"

access(all) fun main(distId: UInt64): {String: String} {
    return PDS.getDistInfo(distId: distId)!.metadata
}
