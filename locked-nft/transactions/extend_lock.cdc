import NFTLocker from "../contracts/NFTLocker.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction(id: UInt64, extendedDuration: UInt64) {
    let lockRef: &NFTLocker.Collection

    prepare(signer: AuthAccount) {
        self.lockRef = signer
            .borrow<&NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
            ?? panic("Account does not store an object at the specified path")
    }

    execute {
            self.lockRef.extendLock(id: id, nftType: Type<@ExampleNFT.NFT>(), extendedDuration: extendedDuration)
    }
}