import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

// Get the UUIDs of the collection added to the parent Aggregator resource
pub fun main(account: Address): [UInt64] {
    let supplierPublicCapability = getAccount(account).getCapability<
        &NFTProviderAggregator.Supplier{NFTProviderAggregator.SupplierPublic}>(
        NFTProviderAggregator.SupplierPublicPath)!

    let supplierPublicRef = supplierPublicCapability.borrow()!

    return supplierPublicRef.getCollectionUUIDs()
}
