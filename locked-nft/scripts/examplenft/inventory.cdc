import NonFungibleToken from "../..contracts/NonFungibleToken.cdc"
import ExampleNFT from "../../contracts/ExampleNFT.cdc"

pub fun main(acctAddress: Address): [UInt64] {
    let nftOwner = getAccount(acctAddress);
    let capability = nftOwner.getCapability<&{NonFungibleToken.CollectionPublic}>(ExampleNFT.CollectionPublicPath);
    let borrowed = capability.borrow() ?? panic("Could not borrow receiver reference")
    return borrowed.getIDs()
}