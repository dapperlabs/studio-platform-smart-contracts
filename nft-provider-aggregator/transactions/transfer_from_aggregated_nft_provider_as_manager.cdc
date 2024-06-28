import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import ExampleNFT from "../contracts/ExampleNFT.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a manager account to transfer a NFT from aggregated NFT provider
/// to the recipient's collection.
///
/// @param recipient: The recipient address of the NFT to transfer.
/// @param withdrawID: The ID of the NFT to withdraw.
///
transaction(recipient: Address, withdrawID: UInt64) {

    let depositRef: &AnyResource{NonFungibleToken.CollectionPublic}
    let aggregatorRef: &NFTProviderAggregator.Aggregator

    prepare(
        manager: AuthAccount,
    ) {
        // Get recipient account
        let recipient = getAccount(recipient)

        // Borrow a public reference to the receivers collection
        self.depositRef = recipient
            .getCapability(ExampleNFT.CollectionPublicPath)!
            .borrow<&AnyResource{NonFungibleToken.CollectionPublic}>()!

        // Create reference to Aggregator
        self.aggregatorRef = manager.borrow<&NFTProviderAggregator.Aggregator>(
            from: NFTProviderAggregator.AggregatorStoragePath)!    
    }

    execute {
        // Withdraw the NFT from the aggregated provider
        let nft <- self.aggregatorRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        self.depositRef.deposit(token: <-nft)
    }
}
