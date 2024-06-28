import NFTProviderAggregator from "NFTProviderAggregator"

// Get the NFT IDs of all the NFTs contained in the collections added to the Aggregator resource
access(all) fun main(account: Address): [UInt64] {
    return getAccount(account).capabilities.borrow<
        &NFTProviderAggregator.Supplier>(
        NFTProviderAggregator.SupplierPublicPath)!.getIDs()
}
