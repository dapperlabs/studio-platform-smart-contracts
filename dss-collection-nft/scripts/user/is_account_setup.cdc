import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import DSSCollection from "../../contracts/DSSCollection.cdc"

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            DSSCollection.DSSCollectionNFTCollectionPublic
        }>(DSSCollection.CollectionPublicPath).check()
}

