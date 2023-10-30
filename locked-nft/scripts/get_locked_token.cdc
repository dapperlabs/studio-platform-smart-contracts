import NFTLocker from "../contracts/NFTLocker.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS

pub fun main(id: UInt64): NFTLocker.LockedData {
    return NFTLocker.getNFTLockerDetails(id: id, nftType: Type<@ExampleNFT.NFT>())!
}