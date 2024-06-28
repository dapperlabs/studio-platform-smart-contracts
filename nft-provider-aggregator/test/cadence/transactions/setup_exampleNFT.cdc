import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import ExampleNFT from "../../contracts/ExampleNFT.cdc"

transaction (
    nftCollectionPublicPathID: String,
    nftCollectionStoragePathID: String
) {

    let storagePath: StoragePath
    let publicPath: PublicPath

    prepare(signer: AuthAccount) {

        self.storagePath = StoragePath(identifier: nftCollectionStoragePathID)
            ?? panic("Provided NFT provider private path has invalid format!")
        self.publicPath = PublicPath(identifier: nftCollectionPublicPathID)
            ?? panic("Provided NFT collection storage path has invalid format!")

        // Return early if the account already has a collection
        if signer.borrow<&ExampleNFT.Collection>(from: self.storagePath) != nil {
            return
        }

        // create a new empty collection
        let collection <- ExampleNFT.createEmptyCollection()

        // save it to the account
        signer.save(<-collection, to: self.storagePath)

        // create a public capability for the collection
        signer.link<&NonFungibleToken.Collection{NonFungibleToken.CollectionPublic}>(self.publicPath, target: self.storagePath)
        assert(signer.getCapability<&{NonFungibleToken.CollectionPublic}>(self.publicPath).check(), message: "did not link pub cap");
    }
}
