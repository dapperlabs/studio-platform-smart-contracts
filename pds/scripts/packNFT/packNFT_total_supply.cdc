import PackNFT from "../../contracts/PackNFT.cdc"

pub fun main(): UInt64{
    return PackNFT.totalSupply 
}
