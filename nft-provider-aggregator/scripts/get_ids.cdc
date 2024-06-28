import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

// Get the NFT IDs of all the NFTs contained in the collections added to the Aggregator resource
pub fun main(account: Address): [UInt64] {
    let supplierPublicCapability = getAccount(account).getCapability<
        &NFTProviderAggregator.Supplier{NFTProviderAggregator.SupplierPublic}>(
        NFTProviderAggregator.SupplierPublicPath)!

    let supplierPublicRef = supplierPublicCapability.borrow()!

    return supplierPublicRef.getIDs()
}
