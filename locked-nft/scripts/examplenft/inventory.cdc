import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(acctAddress: Address): [UInt64] {
    // Get a public collection reference for the owner's account
    let collectionRef = getAccount(acctAddress).capabilities.borrow<
        &ExampleNFT.Collection>(ExampleNFT.CollectionPublicPath)
            ?? panic("Could not borrow a reference of the public collection")

    // Return the NFT IDs in the collection
    return collectionRef.getIDs()
}