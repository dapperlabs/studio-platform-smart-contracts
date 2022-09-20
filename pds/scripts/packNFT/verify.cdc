import PackNFT from "../../contracts/PackNFT.cdc"
import IPackNFT from "../../contracts/IPackNFT.cdc"

pub fun main(id: UInt64, nftString: String): Bool {
    let p = PackNFT.borrowPackRepresentation(id: id) 
    return p!.verify(nftString: nftString)
}
