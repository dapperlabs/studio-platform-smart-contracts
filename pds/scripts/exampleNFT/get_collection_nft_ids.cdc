import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

// This script returns the IDs of the NFTs in the collection
access(all) fun main(address: Address): [UInt64] {
    let account = getAccount(address)

    let collectionRef = getAccount(address).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "cadenceExampleNFTCollection")!)!

    // Return the IDs of the NFTs in the collection
    return collectionRef.getIDs()
}
