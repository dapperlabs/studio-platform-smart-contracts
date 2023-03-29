import NonFungibleToken from "./NonFungibleToken.cdc"
import EnglishPremierLeague from "./EnglishPremierLeague.cdc"

pub fun main(address: Address): Bool {
    let account = getAccount(address)
    return account.getCapability<&{
            NonFungibleToken.CollectionPublic,
            EnglishPremierLeague.MomentNFTCollectionPublic
        }>(EnglishPremierLeague.CollectionPublicPath).check()
}