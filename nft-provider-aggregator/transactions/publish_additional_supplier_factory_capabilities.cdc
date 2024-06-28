import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a manager account to create additional supplier factory capabilities (who are trusted
/// for the ability to use the capability itself but also for potentially copying and saving it somewhere else).
///
/// @param suppliers: Array of the supplier addresses to publish the supplier factory capability to.
/// @param capabilityPublicationIDs: Array of the publication identifiers for each supplier factory capability.
///
transaction(
    suppliers: [Address],
    capabilityPublicationIDs: [String],
    ) {

    prepare(
        manager: AuthAccount,
    ) {
        assert(
            suppliers.length == capabilityPublicationIDs.length,
            message: "suppliers array argument has different size than capabilityIdentifier!"
        )

        // Retrieve supplier factory capability
        let supplierFactoryCapability = manager.getCapability<
        &NFTProviderAggregator.Aggregator{NFTProviderAggregator.SupplierFactory}>(
            NFTProviderAggregator.SupplierFactoryPrivatePath)!

        // Publish supplier factory capability to designated recipients
        for index, recipient in suppliers {
            manager.inbox.publish(
                supplierFactoryCapability,
                name: capabilityPublicationIDs[index],
                recipient: recipient
                )  
        }
    }
}
