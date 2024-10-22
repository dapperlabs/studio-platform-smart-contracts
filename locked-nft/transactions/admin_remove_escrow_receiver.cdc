import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"

/// This transaction removes an escrow leaderboard deposit handler from the admin's ReceiverCollector resource.
///
transaction() {
    // Authorized reference to the NFTLocker ReceiverCollector resource
    let receiverCollectorRef: auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector

    // Name of the receiver to be removed from the ReceiverCollector
    let receiverName: String

    prepare(admin: auth(BorrowValue) &Account) {
        // Set receiver name
        self.receiverName = "add-entry-to-escrow-leaderboard"

        // Borrow an authorized reference to the admin's ReceiverCollector resource
        self.receiverCollectorRef = admin.storage
            .borrow<auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector>(from: NFTLocker.getReceiverCollectorStoragePath())
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        // Add a receiver to the ReceiverCollector with the provided namw, deposit handler, and accepted NFT types
        self.receiverCollectorRef.removeReceiver(name: self.receiverName)
    }

    post {
        self.receiverCollectorRef.getReceiver(name: self.receiverName) == nil : "Receiver not removed"
    }
}
