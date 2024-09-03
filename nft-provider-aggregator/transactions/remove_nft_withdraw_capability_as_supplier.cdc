import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a supplier to remove from the parent Aggregator a NFT provider
/// capability previously added by the same supplier.
///
/// @param collectionUUID: The UUID of the collection to remove the NFT provider capability for.
///
transaction(collectionUUID: UInt64) {

    let supplierRef: auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Supplier

    prepare(
        supplier: auth(BorrowValue) &Account,
    ) {
        // Create reference to Supplier from storage
        self.supplierRef = supplier.storage.borrow<auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath)!
    }

    execute {
        // Remove NFT provider capability by collection UUID
        self.supplierRef.removeNFTWithdrawCapability(collectionUUID: collectionUUID)
    }

    post {
        !self.supplierRef.getCollectionUUIDs().contains(collectionUUID): "NFT Provider capability was not removed!"
    }
}
