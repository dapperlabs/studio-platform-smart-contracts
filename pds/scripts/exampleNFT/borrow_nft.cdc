import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

// This script borrows an NFT from a collection
access(all) fun main(address: Address, id: UInt64) {
    let account = getAccount(address)

    let collectionRef = getAccount(address).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "exampleNFTCollection")!)!

    // Borrow a reference to a specific NFT in the collection
    let _ = collectionRef.borrowNFT(id)
}
