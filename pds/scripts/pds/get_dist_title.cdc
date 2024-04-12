import PDS from "PDS"

access(all) fun main(distId: UInt64): String {
    return PDS.getDistInfo(distId: distId)!.title
}
