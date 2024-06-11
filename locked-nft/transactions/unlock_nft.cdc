import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"
import ExampleNFT from "ExampleNFT"

transaction(id: UInt64) {
    let unlockRef: &NFTLocker.Collection
    let depositRef: &ExampleNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        self.unlockRef = signer.storage
            .borrow<&NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        self.depositRef = signer.storage
            .borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        self.depositRef.deposit(token: <- self.unlockRef.unlock(id: id, nftType: Type<@ExampleNFT.NFT>()))
    }
}