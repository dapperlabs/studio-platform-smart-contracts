import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a manager account to destroy their Aggregator resource
///
transaction() {
    
    prepare(
        manager: AuthAccount,
    ) {
        // Load and destroy the Aggregator resource
        destroy manager.load<@NFTProviderAggregator.Aggregator>(
            from: NFTProviderAggregator.AggregatorStoragePath
            )
    }
}
