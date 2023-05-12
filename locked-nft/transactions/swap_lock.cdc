import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTLocker from "../contracts/NFTLocker.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(nftLockedID: UInt64, nftToLockID: UInt64, duration: UInt64) {
    let exampleCollectionRef: &ExampleNFT.Collection
    let lockRef: &NFTLocker.Collection

    prepare(signer: AuthAccount) {
        self.exampleCollectionRef = signer
            .borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
        self.lockRef = signer
            .borrow<&NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
            ?? panic("Account does not store an object at the specified path")
    }

    execute {
            self.exampleCollectionRef.deposit(token: <- self.lockRef.unlock(id: nftLockedID, nftType: Type<@ExampleNFT.NFT>()))
            self.lockRef.lock(token: <- self.exampleCollectionRef.withdraw(withdrawID: nftToLockID), duration: duration)
    }
}