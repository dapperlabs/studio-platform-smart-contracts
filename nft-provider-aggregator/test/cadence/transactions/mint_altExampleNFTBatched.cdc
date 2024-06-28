import NonFungibleToken from "NonFungibleToken"
import AltExampleNFT from "AltExampleNFT"
import MetadataViews from "MetadataViews"

transaction(recipient: Address, batchSize: Int) {

    let minter: &AltExampleNFT.NFTMinter
    let collectionPublicPath: PublicPath
    let receiver: &AltExampleNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        self.minter = signer.storage.borrow<&AltExampleNFT.NFTMinter>(from: AltExampleNFT.MinterStoragePath)!

        let collectionData = AltExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")
        self.collectionPublicPath = collectionData.publicPath

        self.receiver = getAccount(recipient).capabilities.borrow<&AltExampleNFT.Collection>(self.collectionPublicPath)!
    }

    execute {
        var i = 0
        while i < batchSize {
            self.receiver.deposit(token: <- self.minter.mintNFT(name: "", description: "", thumbnail: "", royalties: []))
            i = i + 1
        }
    }
}
