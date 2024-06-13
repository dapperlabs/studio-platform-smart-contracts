import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(account: Address): [UInt64] {
    let collectionRef = getAccount(account).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "exampleNFTCollection")!)!

    return collectionRef.getIDs()
}
