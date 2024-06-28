import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a manager account to transfer a NFT from aggregated NFT provider
/// to the recipient's collection.
///
/// @param recipient: The recipient address of the NFT to transfer.
/// @param withdrawID: The ID of the NFT to withdraw.
///
transaction(recipient: Address, withdrawID: UInt64) {

    let depositRef: &{NonFungibleToken.CollectionPublic}
    let aggregatorRef: auth(NonFungibleToken.Withdraw) &NFTProviderAggregator.Aggregator

    prepare(
        manager: auth(BorrowValue) &Account,
    ) {
        // Get recipient account
        let recipient = getAccount(recipient)

        // Borrow a public reference to the receivers collection
        self.depositRef = recipient
            .capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(ExampleNFT.CollectionPublicPath)!

        // Create reference to Aggregator
        self.aggregatorRef = manager.storage.borrow<auth(NonFungibleToken.Withdraw) &NFTProviderAggregator.Aggregator>(
            from: NFTProviderAggregator.AggregatorStoragePath)!
    }

    execute {
        // Withdraw the NFT from the aggregated provider
        let nft <- self.aggregatorRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        self.depositRef.deposit(token: <-nft)
    }
}
