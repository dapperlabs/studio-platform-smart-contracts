import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

// This script returns the number of NFTs in the collection of the given address
access(all) fun main(address: Address): Int {
    let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
        ?? panic("ViewResolver does not resolve NFTCollectionData view")

    let account = getAccount(address)

    let collectionRef = getAccount(address).capabilities.borrow<
        &ExampleNFT.Collection>(collectionData.publicPath)!

    return collectionRef.getIDs().length
}
