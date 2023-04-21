import LockedNFT from 0xf8d6e0586b0a20c7

pub fun main(id: UInt64): LockedNFT.LockedData? {
    log("blah blah blah")
    return LockedNFT.getLockedToken(id: id)
}