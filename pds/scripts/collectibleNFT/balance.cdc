import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(account: Address): Int {
    let collectionRef = getAccount(account).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "cadenceExampleNFTCollection")!)!

    return collectionRef.getIDs().length
}
