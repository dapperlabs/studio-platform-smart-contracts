import NFTProviderAggregator from "NFTProviderAggregator"
import Burner from "Burner"

/// Transaction signed by a supplier account to destroy their Supplier resource
///
transaction() {

    prepare(
        supplier: auth(LoadValue) &Account,
    ) {
        // Load and destroy the Aggregator resource
        Burner.burn(<- supplier.storage.load<@NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath
            )
        )
    }
}