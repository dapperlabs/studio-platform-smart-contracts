import NFTLocker from "../contracts/NFTLocker.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS

transaction(id: UInt64) {
    let adminRef: &NFTLocker.Admin

    prepare(signer: AuthAccount) {
        self.adminRef = signer
            .borrow<&NFTLocker.Admin>(from: NFTLocker.GetAdminStoragePath())
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        self.adminRef.expireLock(id: id, nftType: Type<@ExampleNFT.NFT>())
    }
}