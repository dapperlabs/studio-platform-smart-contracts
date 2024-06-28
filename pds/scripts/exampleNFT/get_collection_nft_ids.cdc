import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

// This script returns the IDs of the NFTs in the collection
access(all) fun main(address: Address): [UInt64] {
    let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
        ?? panic("ViewResolver does not resolve NFTCollectionData view")

    let account = getAccount(address)

    let collectionRef = getAccount(address).capabilities.borrow<
        &ExampleNFT.Collection>(collectionData.publicPath)!

    // Return the IDs of the NFTs in the collection
    return collectionRef.getIDs()
}
