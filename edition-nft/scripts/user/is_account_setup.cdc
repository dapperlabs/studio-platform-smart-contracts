import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import EditionNFT from "../../contracts/EditionNFT.cdc"

// Check to see if an account looks like it has been set up to hold EditionNFTs.

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            EditionNFT.EditionNFTCollectionPublic
        }>(EditionNFT.CollectionPublicPath)
        != nil
}

