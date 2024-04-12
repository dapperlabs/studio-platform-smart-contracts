import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"

access(all) fun main(id: UInt64): UInt8 {
    let p = PackNFT.borrowPackRepresentation(id: id)
    return p!.status.rawValue
}
