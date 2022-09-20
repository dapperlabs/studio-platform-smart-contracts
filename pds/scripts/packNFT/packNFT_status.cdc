import PackNFT from "../../contracts/PackNFT.cdc"
import IPackNFT from "../../contracts/IPackNFT.cdc"

pub fun main(id: UInt64): UInt8 {
    let p = PackNFT.borrowPackRepresentation(id: id) 
    return p!.status.rawValue
}
