import NFTLocker from "../contracts/NFTLocker.cdc"

pub fun main(id: UInt64, nftType: Type): NFTLocker.LockedData? {
    return NFTLocker.getNFTLockerDetails(id: id, nftType: nftType)
}