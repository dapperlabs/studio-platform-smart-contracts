import NFTProviderAggregator from "NFTProviderAggregator"

// Get the UUIDs of the collection added to the parent Aggregator resource
access(all) fun main(account: Address): [UInt64] {
    return getAuthAccount<auth(BorrowValue) &Account>(account).storage.borrow<
        &NFTProviderAggregator.Aggregator>(from:
        NFTProviderAggregator.AggregatorStoragePath)!.getCollectionUUIDs()
}
