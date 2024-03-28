import NonFungibleToken from "NonFungibleToken"
import PackNFT from "PackNFT"

transaction(recipient: Address, withdrawID: UInt64) {
    // Reference to the withdrawer's collection
    let withdrawRef: auth(NonFungibleToken.Withdraw) &PackNFT.Collection

    // Reference of the collection to deposit the NFT to
    let receiverRef: &PackNFT.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        let recipient = getAccount(recipient)

        // borrow a reference to the signer's NFT collection
        self.withdrawRef = signer.storage.borrow<
            auth(NonFungibleToken.Withdraw) &PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!

        // borrow a public reference to the receivers collection
        self.receiverRef = recipient
            .capabilities.borrow<&PackNFT.Collection>(PackNFT.CollectionPublicPath)!
    }

    execute {
        // Withdraw the NFT from the owner's collection and deposit it in the recipient's collection
        self.receiverRef.deposit(token: <- self.withdrawRef.withdraw(withdrawID: withdrawID))
    }
}
