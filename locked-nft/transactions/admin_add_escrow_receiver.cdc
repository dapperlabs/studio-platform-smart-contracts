import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"
import Escrow from "Escrow"
import ExampleNFT from "ExampleNFT"

/// This transaction creates a new ReceiverCollector resource and adds a new receiver to it with a deposit method that adds an NFT to an escrow leaderboard.
///
transaction() {
    // Authorized reference to the NFTLocker ReceiverCollector resource
    let receiverCollectorRef: auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector

    prepare(admin: auth(SaveValue, BorrowValue) &Account) {
        // Check if the ReceiverCollector resource does not exist
        if NFTLocker.borrowAdminReceiverCollectorPublic() == nil {
            // Borrow a reference to the NFTLocker Admin resource
            let adminRef = admin.storage.borrow<&NFTLocker.Admin>(from: NFTLocker.getAdminStoragePath())
                ?? panic("Could not borrow a reference to the owner's collection")

            // Create a new ReceiverCollector resource and save it in storage
            admin.storage.save(<- adminRef.createReceiverCollector(), to: NFTLocker.getReceiverCollectorStoragePath())
        }

        // Borrow an authorized reference to the NFTLocker ReceiverCollector resource
        self.receiverCollectorRef = admin.storage
            .borrow<auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector>(from: NFTLocker.getReceiverCollectorStoragePath())
            ?? panic("Could not borrow a reference to the owner's collection")
    }

    execute {
        // Add a new receiver to the ReceiverCollector with the provided deposit wrapper and accepted NFT types
        self.receiverCollectorRef.addReceiver(
            name: "add-entry-to-escrow-leaderboard",
            authorizedDepositHandler: Escrow.DepositHandler(),
            eligibleNFTTypes: {Type<@ExampleNFT.NFT>(): true}
        )
    }
}
