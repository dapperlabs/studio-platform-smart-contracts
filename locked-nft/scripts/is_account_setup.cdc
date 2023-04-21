import NonFungibleToken from "./NonFungibleToken.cdc"
import LockedNFT from "./LockedNFT.cdc"

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            LockedNFT.LockedNFTCollectionPublic
        }>(LockedNFT.CollectionPublicPath).check()
}