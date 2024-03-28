import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"

access(all) fun main(id: UInt64, nftString: String): Bool {
    let p = PackNFT.borrowPackRepresentation(id: id)
    return p!.verify(nftString: nftString)
}
