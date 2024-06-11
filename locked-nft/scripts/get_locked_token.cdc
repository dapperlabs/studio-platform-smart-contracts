import NFTLocker from "NFTLocker"
import ExampleNFT from "ExampleNFT"

access(all) fun main(id: UInt64): NFTLocker.LockedData? {
    return NFTLocker.getNFTLockerDetails(id: id, nftType: Type<@ExampleNFT.NFT>())
}