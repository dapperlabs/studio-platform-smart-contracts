import NFTProviderAggregator from "NFTProviderAggregator"

// Get the UUIDs of the collection added by the provided account
access(all) fun main(account: Address): [UInt64] {
        return getAccount(account).capabilities.borrow<
        &NFTProviderAggregator.Supplier>(
        NFTProviderAggregator.SupplierPublicPath)!.getSupplierAddedCollectionUUIDs()
}
