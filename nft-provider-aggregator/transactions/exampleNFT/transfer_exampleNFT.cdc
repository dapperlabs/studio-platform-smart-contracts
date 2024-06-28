import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

transaction(recipient: Address, withdrawID: UInt64) {
    // Reference to the withdrawer's collection
    let withdrawRef: auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection

    // Reference of the collection to deposit the NFT to
    let receiverRef: &ExampleNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        let recipient = getAccount(recipient)

        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        // borrow a reference to the signer's NFT collection
        self.withdrawRef = signer.storage.borrow<
            auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection>(from: collectionData.storagePath)!

        // borrow a public reference to the receivers collection
        self.receiverRef = recipient
            .capabilities.borrow<&ExampleNFT.Collection>(collectionData.publicPath)!
    }

    execute {
        // Withdraw the NFT from the owner's collection and deposit it in the recipient's collection
        self.receiverRef.deposit(token: <- self.withdrawRef.withdraw(withdrawID: withdrawID))
    }
}
