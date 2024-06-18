import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"
import ExampleNFT from "ExampleNFT"

transaction(id: UInt64) {
    let unlockRef: auth(NFTLocker.Operate) &NFTLocker.Collection
    let depositRef: &ExampleNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow an authorized reference to the owner's NFTLocker collection
        self.unlockRef = signer.storage
            .borrow<auth(NFTLocker.Operate) &NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        // Borrow a reference to the owner's ExampleNFT collection
        self.depositRef = signer.storage
            .borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        self.depositRef.deposit(token: <- self.unlockRef.unlock(id: id, nftType: Type<@ExampleNFT.NFT>()))
    }
}