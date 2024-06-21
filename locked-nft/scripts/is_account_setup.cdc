import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"

access(all) fun main(address: Address): Bool {
    return getAccount(address).capabilities.get<&{NFTLocker.LockedCollection}>(NFTLocker.CollectionPublicPath).check()
}