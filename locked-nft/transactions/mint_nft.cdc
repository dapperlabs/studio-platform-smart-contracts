import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

transaction(recipient: Address) {
    let minter: &ExampleNFT.NFTMinter
    let recipientCollectionRef: &{NonFungibleToken.CollectionPublic}

    prepare(signer: auth(BorrowValue) &Account) {

        self.minter = signer.storage.borrow<&ExampleNFT.NFTMinter>(from: ExampleNFT.MinterStoragePath)
            ?? panic("Account does not store an object at the specified path")

        self.recipientCollectionRef = getAccount(recipient)
            .capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(ExampleNFT.CollectionPublicPath)
            ?? panic("Could not get receiver reference to the NFT Collection")
    }

    execute {
        // Stub data for other parameters
        let name: String = "Example NFT"
        let description: String = "This is an example NFT."
        let thumbnail: String = "nft.jpg"

        self.recipientCollectionRef.deposit(token: <- self.minter.mintNFT(
            name: name,
            description: description,
            thumbnail: thumbnail,
            royalties: []
        )
        )
    }

    // post {
    //     self.recipientCollectionRef.getIDs().contains(self.mintingIDBefore): "The next NFT ID should have been minted and delivered"
    //     ExampleNFT.totalSupply == self.mintingIDBefore + 1: "The total supply should have been increased by 1"
    // }
}