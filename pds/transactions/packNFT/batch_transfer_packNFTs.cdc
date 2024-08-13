import NonFungibleToken from "NonFungibleToken"
import PackNFT from "PackNFT"
import DapperStorageRent from "DapperStorageRent"

// This transaction transfers NFL NFTs from one account to another.

transaction(recipients: [Address], withdrawIDs: [UInt64]) {
    // Reference to the withdrawer's collection
    let withdrawRef: auth(NonFungibleToken.Withdraw) &PackNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // borrow a reference to the signer's NFT collection
        self.withdrawRef = signer.storage.borrow<
            auth(NonFungibleToken.Withdraw) &PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
    }

    execute {
        var i = 0
        while i < withdrawIDs.length {
            // borrow a public reference to the receiver collection
            let depositRef = getAccount(recipients[i])
                .capabilities.borrow<&PackNFT.Collection>(PackNFT.CollectionPublicPath)!

            // try to refill the receiver's storage
            DapperStorageRent.tryRefill(recipients[i])

            depositRef.deposit(token: <- self.withdrawRef.withdraw(withdrawID: withdrawIDs[i]))
            i = i + 1
        }
    }
}
