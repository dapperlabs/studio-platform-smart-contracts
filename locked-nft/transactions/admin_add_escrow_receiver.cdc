import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"
import Escrow from "Escrow"
import ExampleNFT from "ExampleNFT"

/// This transaction creates a ReceiverCollector resource and adds an escrow leaderboard deposit handler to it.
///
transaction() {
    // Authorized reference to the NFTLocker ReceiverCollector resource
    let receiverCollectorRef: auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector

    // Name of the receiver to be added to the ReceiverCollector
    let receiverName: String

    prepare(admin: auth(SaveValue, BorrowValue) &Account) {
        // Set receiver name
        self.receiverName = "add-entry-to-escrow-leaderboard"

        // Create a ReceiverCollector resource if it does not exist yet in admin storage
        if NFTLocker.borrowAdminReceiverCollectorPublic() == nil {
            // Borrow a reference to the NFTLocker Admin resource
            let adminRef = admin.storage.borrow<&NFTLocker.Admin>(from: NFTLocker.GetAdminStoragePath())
                ?? panic("Could not borrow a reference to the owner's collection")

            // Create a ReceiverCollector resource and save it in storage
            admin.storage.save(<- adminRef.createReceiverCollector(), to: NFTLocker.getReceiverCollectorStoragePath())
        }

        // Borrow an authorized reference to the admin's ReceiverCollector resource
        self.receiverCollectorRef = admin.storage
            .borrow<auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector>(from: NFTLocker.getReceiverCollectorStoragePath())
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        // Add a receiver to the ReceiverCollector with the provided namw, deposit handler, and accepted NFT types
        self.receiverCollectorRef.addReceiver(
            name: self.receiverName,
            authorizedDepositHandler: Escrow.DepositHandler(),
            eligibleNFTTypes: {Type<@ExampleNFT.NFT>(): true}
        )
    }

    post {
        self.receiverCollectorRef.getReceiver(name: self.receiverName) != nil : "Receiver not added"
    }
}
