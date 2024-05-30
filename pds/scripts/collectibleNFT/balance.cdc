import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(account: Address): Int {
    let collectionRef = getAccount(account).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "exampleNFTCollection")!)!

    return collectionRef.getIDs().length
}
