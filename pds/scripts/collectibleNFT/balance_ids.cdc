import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"

access(all) fun main(account: Address, offset: UInt64, limit: UInt64): [UInt64] {
    let collectionRef = getAccount(account).capabilities.borrow<
        &ExampleNFT.Collection>(PublicPath(identifier: "cadenceExampleNFTCollection")!)!

    let ids = collectionRef.getIDs()
    let idsLen = UInt64(ids.length)

    var res: [UInt64] = []
    var i = offset
    while i < offset+limit && i < idsLen {
        res.append(ids[i])
        i = i + 1
    }

    return res
}
