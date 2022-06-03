import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Golazo from "../../contracts/Golazo.cdc"

// Check to see if an account looks like it has been set up to hold Golazo NFTs.

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            Golazo.MomentNFTCollectionPublic
        }>(Golazo.CollectionPublicPath)
        != nil
}

