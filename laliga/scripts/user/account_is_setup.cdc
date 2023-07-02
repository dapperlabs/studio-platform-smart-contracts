import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import Golazos from "../../contracts/Golazos.cdc"
import MetadataViews from 0x{{.MetadataViewsAddress}}

// Check to see if an account looks like it has been set up to hold Golazos NFTs.

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            Golazos.MomentNFTCollectionPublic,
            MetadataViews.ResolverCollection
        }>(Golazos.CollectionPublicPath)
        != nil
}

