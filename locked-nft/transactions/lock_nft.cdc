import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTLocker from "../contracts/NFTLocker.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(nftID: UInt64, duration: UInt64) {
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
            self.lockRef.lock(token: <- self.exampleCollectionRef.withdraw(withdrawID: nftID), duration: duration)
    }
}