import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a supplier to remove from the parent Aggregator a NFT provider
/// capability previously added by the same supplier.
///
/// @param collectionUUID: The UUID of the collection to remove the NFT provider capability for.
///
transaction(collectionUUID: UInt64) {

    let supplierRef: &NFTProviderAggregator.Supplier

    prepare(
        supplier: AuthAccount,
    ) {
        // Create reference to Supplier from storage
        self.supplierRef = supplier.borrow<&NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath)!
    }

    execute {
        // Remove NFT provider capability by collection UUID
        self.supplierRef.removeNFTProviderCapability(collectionUUID: collectionUUID)
    }

    post {
        !self.supplierRef.getCollectionUUIDs().contains(collectionUUID): "NFT Provider capability was not removed!"
    }
}
