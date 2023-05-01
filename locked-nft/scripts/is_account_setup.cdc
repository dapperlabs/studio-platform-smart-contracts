import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTLocker from "../contracts/NFTLocker.cdc"

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NFTLocker.LockedCollection
        }>(NFTLocker.CollectionPublicPath).check()
}