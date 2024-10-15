import NFTLocker from "NFTLocker"
import ExampleNFT from "ExampleNFT"

transaction(id: UInt64) {
    let adminRef: &NFTLocker.Admin

    prepare(signer: auth(BorrowValue) &Account) {
        self.adminRef = signer.storage
            .borrow<&NFTLocker.Admin>(from: NFTLocker.getAdminStoragePath())
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        self.adminRef.expireLock(id: id, nftType: Type<@ExampleNFT.NFT>())
    }
}