import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a supplier account to destroy their Supplier resource
///
transaction() {
    
    prepare(
        supplier: AuthAccount,
    ) {
        // Load and destroy the Aggregator resource
        destroy supplier.load<@NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath
            )
    }
}