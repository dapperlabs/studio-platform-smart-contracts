import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(account: Address): Int {
    let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
        ?? panic("ViewResolver does not resolve NFTCollectionData view")

    let collectionRef = getAccount(account).capabilities.borrow<
        &ExampleNFT.Collection>(collectionData.publicPath)!

    return collectionRef.getIDs().length
}
