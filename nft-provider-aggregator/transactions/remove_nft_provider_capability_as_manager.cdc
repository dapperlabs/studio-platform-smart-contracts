import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a manager to remove a NFT provider capability from the Aggregator resource.
///
/// @param collectionUUID: The UUID of the collection to remove the NFT provider capability for.
///
transaction(collectionUUID: UInt64) {

    let aggregatorRef: auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Aggregator

    prepare(
        manager: auth(BorrowValue) &Account,
    ) {
        // Create reference to Aggregator from storage
        self.aggregatorRef = manager.storage.borrow<auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Aggregator>(
            from: NFTProviderAggregator.AggregatorStoragePath)!
    }

    execute {
        // Remove NFT provider capability by collection UUID
        self.aggregatorRef.removeNFTWithdrawCapability(collectionUUID: collectionUUID)
    }

    post {
        !self.aggregatorRef.getCollectionUUIDs().contains(collectionUUID): "NFT Provider capability was not removed!"
    }
}
