import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by any account to create a new Aggregator resource, save it
/// in the account's storage, and publish a supplier factory capability to each designated
/// supplier (who are trusted for the ability to use the capability itself but also for potentially
/// copying and saving it somewhere else).
///
/// @param nftTypeIdentifier: The identifier of the NFT collection type to be used by the Aggregator resource
/// following the format "A.<account-address>.<NFT-contract-name>.Collection" (e.g., "A.01cf0e2f2f715450.ExampleNFT.Collection").
/// @param suppliers: Array of the supplier addresses to publish the supplier factory capability to.
/// @param capabilityPublicationIDs: Array of the publication identifiers for each supplier factory capability.
///
transaction(
    nftTypeIdentifier: String,
    suppliers: [Address],
    capabilityPublicationIDs: [String],
    ) {

    prepare(
        signer: auth(BorrowValue, SaveValue, Capabilities, Inbox) &Account,
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
        let supplierAccessCapability = signer.capabilities.storage.issue<auth(NFTProviderAggregator.Operate) &{NFTProviderAggregator.SupplierAccess}>(
            NFTProviderAggregator.AggregatorStoragePath
        )

        // Save Aggregator resource capability to storage
        signer.storage.save(
            supplierAccessCapability,
            to: StoragePath(identifier: NFTProviderAggregator.SupplierStoragePath.toString().replaceAll(of: "storage/", with: "").concat("_PrivateCap"))!
        )

        // Create Aggregator resource and save to storage
        let aggregator <- NFTProviderAggregator.createAggregator(
            nftTypeIdentifier: nftTypeIdentifier,
            supplierAccessCapability: supplierAccessCapability
            )
        signer.storage.save(<-aggregator, to: NFTProviderAggregator.AggregatorStoragePath)

        // Create supplier factory capability
        let supplierFactoryCapability = signer.capabilities.storage.issue<auth(NFTProviderAggregator.Operate) &{NFTProviderAggregator.SupplierFactory}>(
            NFTProviderAggregator.AggregatorStoragePath
        )

        // Save supplier resource capability to storage
        signer.storage.save(
            supplierFactoryCapability,
            to: StoragePath(identifier: NFTProviderAggregator.AggregatorStoragePath.toString().replaceAll(of: "storage/", with: "").concat("_PrivateCap"))!
        )

        // Publish supplier factory capability to designated recipients
        for index, recipient in suppliers {
            signer.inbox.publish(
                supplierFactoryCapability,
                name: capabilityPublicationIDs[index],
                recipient: recipient)
        }
    }
}
