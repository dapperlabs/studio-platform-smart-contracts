import NFTProviderAggregator from "NFTProviderAggregator"

// Get the UUIDs of the parent Aggregator resource
access(all) fun main(account: Address): UInt64 {
    return getAccount(account).capabilities.borrow<
        &NFTProviderAggregator.Supplier>(
        NFTProviderAggregator.SupplierPublicPath)!.getAggregatorUUID()
}
