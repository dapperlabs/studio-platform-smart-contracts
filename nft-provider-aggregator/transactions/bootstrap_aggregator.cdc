import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by any account to create a new Aggregator resource, save it
/// in the account's storage, and publish a supplier factory capability to each designated
/// supplier (who are trusted for the ability to use the capability itself but also for potentially
/// copying and saving it somewhere else).
/// 
/// @param nftTypeIdentifier: The identifier of the NFT collection type to be used by the Aggregator resource
/// following the format "A.<account-address>.<NFT-contract-name>.Collection" (e.g., "A.01cf0e2f2f715450.ExampleNFT.Collection").
/// @param suppliers: Array of the supplier addresses to publish the supplier factory capability to.
/// @param capabilityPublicationIDs: Array of the publication identifiers for each supplier factory capability.
/// @param useBorrowNFTSafe: Boolean to choose whether or not to use the borrowNFTSafe function when checking if an NFT exists
/// (using it may be useful for large collections where calling getIDs fails due to exceeding the gas limit).
///
transaction(
    nftTypeIdentifier: String,
    suppliers: [Address],
    capabilityPublicationIDs: [String],
    useBorrowNFTSafe: Bool
    ) {

    prepare(
        signer: AuthAccount,
    ) {
        assert(
            suppliers.length == capabilityPublicationIDs.length,
            message: "suppliers array argument has different size than capabilityIdentifier!"
        )
        assert(
            nftTypeIdentifier.slice(from:0, upTo:2) == "A.",
            message: "Invalid nftTypeIdentifier format. Must follow the format: 'A.<account-address>.<NFT-contract-name>.Collection'."
        )

        // Create supplier access capability (used for Supplier resource factory)
        let supplierAccessCapability = signer.link<&NFTProviderAggregator.Aggregator{NFTProviderAggregator.SupplierAccess}>(
            NFTProviderAggregator.SupplierAccessPrivatePath,
            target: NFTProviderAggregator.AggregatorStoragePath
        ) ?? panic("Could not link Aggregator capability!")

        // Create Aggregator resource and save to storage
        let aggregator <- NFTProviderAggregator.createAggregator(
            nftTypeIdentifier: nftTypeIdentifier,
            useBorrowNFTSafe: useBorrowNFTSafe,
            supplierAccessCapability: supplierAccessCapability
            )
        signer.save(<-aggregator, to: NFTProviderAggregator.AggregatorStoragePath)

        // Create supplier factory capability
        let supplierFactoryCapability = signer.link<&NFTProviderAggregator.Aggregator{NFTProviderAggregator.SupplierFactory}>(
            NFTProviderAggregator.SupplierFactoryPrivatePath,
            target: NFTProviderAggregator.AggregatorStoragePath
        ) ?? panic("Could not link Aggregator capability!")

        // Publish supplier factory capability to designated recipients
        for index, recipient in suppliers {
            signer.inbox.publish(
                supplierFactoryCapability,
                name: capabilityPublicationIDs[index],
                recipient: recipient)  
        }
    }
}
