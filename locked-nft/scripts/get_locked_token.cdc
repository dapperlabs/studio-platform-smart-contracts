import LockedNFT from "../contracts/LockedNFT.cdc"

pub fun main(id: UInt64): LockedNFT.LockedData? {
    return LockedNFT.getLockedToken(id: id)
}