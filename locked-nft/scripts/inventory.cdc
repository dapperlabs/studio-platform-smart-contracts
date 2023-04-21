import LockedNFT from "../contracts/LockedNFT.cdc"

pub fun main(acctAddress: Address): [UInt64] {
    let nftOwner = getAccount(acctAddress);
    let capability = nftOwner.getCapability<&{LockedNFT.LockedCollection}>(LockedNFT.CollectionPublicPath);
    let borrowed = capability.borrow() ?? panic("Could not borrow receiver reference")
    return borrowed.getIDs()
}