import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTLocker from "../contracts/NFTLocker.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS

transaction(id: UInt64) {
    let unlockRef: &NFTLocker.Collection
    let depositRef: &ExampleNFT.Collection

    prepare(signer: AuthAccount) {
        self.unlockRef = signer
            .borrow<&NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        self.depositRef = signer
            .borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        self.depositRef.deposit(token: <- self.unlockRef.unlock(id: id, nftType: Type<@ExampleNFT.NFT>()))
    }
}