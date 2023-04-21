import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import LockedNFT from "../contracts/LockedNFT.cdc"

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            LockedNFT.LockedCollection
        }>(LockedNFT.CollectionPublicPath).check()
}