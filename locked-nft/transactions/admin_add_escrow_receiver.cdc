import NonFungibleToken from "NonFungibleToken"
import NFTLocker from "NFTLocker"
import Escrow from "Escrow"
import ExampleNFT from "ExampleNFT"

/// This transaction creates a new ReceiverCollector resource and adds a new receiver to it with a deposit method that adds an NFT to an escrow leaderboard.
///
transaction() {
    // Auhtorized reference to the NFTLocker ReceiverCollector resource
    let receiverCollectorRef: auth(NFTLocker.Operate) &NFTLocker.ReceiverCollector

    // Deposit method to be added to the ReceiverCollector resource
    let depositMethod: fun(@{NonFungibleToken.NFT}, NFTLocker.LockedData, {String: AnyStruct})

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

        // Define the deposit method to be used by the Receiver
        self.depositMethod = fun(nft: @{NonFungibleToken.NFT}, lockedTokenDetails: NFTLocker.LockedData, passThruParams: {String: AnyStruct}) {
            // Get leaderboard name from pass-thru parameters
            let leaderboardName = passThruParams["leaderboardName"] as? String
                ?? panic("Missing or invalid leaderboard name")

            // Get the Escrow contract account
            let escrowAccount = getAccount(Address.fromString(Type<Escrow>().identifier.slice(from: 2, upTo: 18))!)

            // Get the Escrow Collection public reference
            let escrowCollectionPublic = escrowAccount.capabilities.borrow<&Escrow.Collection>(Escrow.CollectionPublicPath)
                ?? panic("Could not borrow a reference to the public leaderboard collection")

            // Add the NFT to the escrow leaderboard
            escrowCollectionPublic.addEntryToLeaderboard(nft: <-nft, leaderboardName: leaderboardName, ownerAddress: lockedTokenDetails.owner)
        }
    }

    execute {
        // Add a new receiver to the ReceiverCollector with the provided deposit method and accepted NFT types
        self.receiverCollectorRef.addReceiver(
            name: "add-entry-to-escrow-leaderboard",
            depositMethod: self.depositMethod,
            eligibleNFTTypes: {Type<@ExampleNFT.NFT>(): true}
        )
    }
}