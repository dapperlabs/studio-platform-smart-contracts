import NFTProviderAggregator from "NFTProviderAggregator"
import Burner from "Burner"

/// Transaction signed by a supplier account to destroy their Supplier resource
///
transaction() {

    let supplierResource: @NFTProviderAggregator.Supplier
    let supplierAddedCollectionUUIDsCount: Int
    let aggregatorRef: &NFTProviderAggregator.Aggregator
    let aggregatorCollectionUUIDsCountBefore: Int

    prepare(
        supplier: auth(LoadValue) &Account,
    ) {
        // Load the Supplier resource
        self.supplierResource <- supplier.storage.load<@NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath
            )
            ?? panic("Supplier does not exist")

        // Get aggregator and supplier collection UUIDs counts
        self.aggregatorRef = self.supplierResource.borrowPublicAggregator()
        self.supplierAddedCollectionUUIDsCount = self.supplierResource.getSupplierAddedCollectionUUIDs().length
        self.aggregatorCollectionUUIDsCountBefore = self.aggregatorRef.getCollectionUUIDs().length
    }

    execute {
        // Destroy the Supplier resource
        let burnableResource <- self.supplierResource as @{Burner.Burnable}
        Burner.burn(<- burnableResource)
    }

    post {
        self.aggregatorRef.getCollectionUUIDs().length == self.aggregatorCollectionUUIDsCountBefore + self.supplierAddedCollectionUUIDsCount:
            "Supplier collection providers were not removed from the Aggregator"
    }
}