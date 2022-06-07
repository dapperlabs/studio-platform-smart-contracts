import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Sport from "../../contracts/Sport.cdc"

// Check to see if an account looks like it has been set up to hold Sport NFTs.

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            Sport.MomentNFTCollectionPublic
        }>(Sport.CollectionPublicPath)
        != nil
}

