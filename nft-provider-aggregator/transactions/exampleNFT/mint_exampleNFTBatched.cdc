import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

transaction(recipient: Address, batchSize: Int) {

    let minter: &ExampleNFT.NFTMinter
    let collectionPublicPath: PublicPath
    let receiver: &ExampleNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        self.minter = signer.storage.borrow<&ExampleNFT.NFTMinter>(from: ExampleNFT.MinterStoragePath)!

        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")
        self.collectionPublicPath = collectionData.publicPath

        self.receiver = getAccount(recipient).capabilities.borrow<&ExampleNFT.Collection>(self.collectionPublicPath)!
    }

    execute {
        var i = 0
        while i < batchSize {
            self.receiver.deposit(token: <- self.minter.mintNFT(name: "", description: "", thumbnail: "", royalties: []))
            i = i + 1
        }
    }
}
