import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

// This script returns the number of NFTs in the collection of the given address
access(all) fun main(address: Address): Int {
    let account = getAccount(address)

    let collectionRef = getAccount(address).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "cadenceExampleNFTCollection")!)!

    return collectionRef.getIDs().length
}
