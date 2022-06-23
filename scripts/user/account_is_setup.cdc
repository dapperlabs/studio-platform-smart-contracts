import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DapperSport from "../../contracts/DapperSport.cdc"
import MetadataViews from 0xMETADATAVIEWSADDRESS

// Check to see if an account looks like it has been set up to hold DapperSport NFTs.

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            DapperSport.MomentNFTCollectionPublic,
            MetadataViews.ResolverCollection
        }>(DapperSport.CollectionPublicPath)
        != nil
}

