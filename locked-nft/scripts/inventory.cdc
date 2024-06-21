import NFTLocker from "NFTLocker"
import ExampleNFT from "ExampleNFT"

access(all) fun main(acctAddress: Address): [UInt64]? {
    let lockedCollectionRef = getAccount(acctAddress).capabilities.borrow<&{NFTLocker.LockedCollection}>(NFTLocker.CollectionPublicPath)
        ?? panic("Could not borrow receiver reference")
    return lockedCollectionRef.getIDs(nftType: Type<@ExampleNFT.NFT>())
}