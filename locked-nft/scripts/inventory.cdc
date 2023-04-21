import LockedNFT from 0xf8d6e0586b0a20c7

pub fun main(acctAddress: Address): [UInt64] {
    let nftOwner = getAccount(acctAddress);
    let capability = nftOwner.getCapability<&{LockedNFT.LockedCollection}>(LockedNFT.CollectionPublicPath);
    let borrowed = capability.borrow() ?? panic("Could not borrow receiver reference")
    return borrowed.getIDs()
}