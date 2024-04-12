import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"


// Used for testing purposes

transaction(recipient: Address, batchSize: Int) {

    /// local variable for storing the minter reference
    let minter: &ExampleNFT.NFTMinter

    /// Reference to the receiver's collection
    let recipientCollectionRef: &{NonFungibleToken.Receiver}

    prepare(signer: auth(BorrowValue) &Account) {
        let collectionData = ExampleNFT.resolveContractView(
            resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        // borrow a reference to the NFTMinter resource in storage
        self.minter = signer.storage.borrow<&ExampleNFT.NFTMinter>(from: ExampleNFT.MinterStoragePath)
            ?? panic("Account does not store an object at the specified path")

        // Borrow the recipient's public NFT collection reference
        self.recipientCollectionRef = getAccount(recipient).capabilities.borrow<&{NonFungibleToken.Receiver}>(
                collectionData.publicPath
            ) ?? panic("Could not get receiver reference to the NFT Collection")
    }

    execute {
        var i = 0
        while i < batchSize {
            self.recipientCollectionRef.deposit(
                token: <-self.minter.mintNFT(
                    name: "",
                    description: "",
                    thumbnail: "",
                    royalties: []
        )
            )
            i = i + 1
        }
    }
}
