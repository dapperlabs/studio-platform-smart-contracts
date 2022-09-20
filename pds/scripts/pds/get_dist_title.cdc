import PDS from "../../contracts/PDS.cdc"

pub fun main(distId: UInt64): String {
    return PDS.getDistInfo(distId: distId)!.title
}
