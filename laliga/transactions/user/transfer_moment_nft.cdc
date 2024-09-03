import NonFungibleToken from "NonFungibleToken"
import Golazos from "Golazos"

/// This transaction transfers the Golazos NFT with the given ID from the signer's collection
/// to the recipient's collection.
///
transaction(recipientAddress: Address, withdrawID: UInt64) {

    // Reference to the withdrawer's collection
    let withdrawRef: auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}

    // Reference of the collection to deposit the NFT to
    let receiverRef: &Golazos.Collection

    prepare(signer: auth(BorrowValue) &Account) {
        // Borrow a reference to the signer's NFT collection
        self.withdrawRef = signer.storage.borrow<auth(NonFungibleToken.Withdraw)
            &{NonFungibleToken.Collection}>(from: Golazos.CollectionStoragePath)
                ?? panic("Could not borrow a reference to the owner's collection")

        // Borrow a public reference to the receivers collection
        self.receiverRef = getAccount(recipientAddress).capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
            ?? panic("Could not borrow a reference to the collection receiver")
    }

    execute {
        // Withdraw the NFT from the owner's collection and deposit it in the recipient's collection
        self.receiverRef.deposit(token: <- self.withdrawRef.withdraw(withdrawID: withdrawID))
    }
}