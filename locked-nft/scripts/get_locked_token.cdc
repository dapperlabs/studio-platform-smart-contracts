import NFTLocker from "../contracts/NFTLocker.cdc"

pub fun main(id: UInt64): NFTLocker.LockedData? {
    return NFTLocker.getLockedToken(id: id)
}