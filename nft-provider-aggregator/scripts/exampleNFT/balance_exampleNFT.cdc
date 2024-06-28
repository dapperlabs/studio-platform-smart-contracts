import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(account: Address): [UInt64] {
    return getAccount(account)
        .capabilities.borrow<&ExampleNFT.Collection>(ExampleNFT.CollectionPublicPath)!.getIDs()
}
