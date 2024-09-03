import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"
import ExampleNFT from "ExampleNFT"


transaction(nftID: UInt64, duration: UInt64) {
    let exampleCollectionRef: auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection
    let lockRef: auth(NFTLocker.Operate) &NFTLocker.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow an authorized reference to the signer's ExplampleNFT collection
        self.exampleCollectionRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw)
            &ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // Borrow an authorized reference to the signer's NFTLocker collection
        self.lockRef = signer.storage.borrow<auth(NFTLocker.Operate)
            &NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
            ?? panic("Account does not store an object at the specified path")
    }

    execute {
            self.lockRef.lock(token: <- self.exampleCollectionRef.withdraw(withdrawID: nftID), duration: duration)
    }
}