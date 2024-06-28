import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a supplier to claim the supplier factory capability (must
/// have been previously published by the manager), create a new Supplier resource, and
/// save it in the supplier's storage.
/// Note: A published capability can only be claimed once.
///
/// @param manager: The address of the manager that published the supplier factory capability.
/// @param capabilityPublicationID: The publication identifier of the supplier factory capability.
///
transaction(
    manager: Address,
    capabilityPublicationID: String,
    ) {
    
    let supplierFactoryRef: &NFTProviderAggregator.Aggregator{NFTProviderAggregator.SupplierFactory}
    let supplierPublicCapability: Capability<
        &NFTProviderAggregator.Supplier{NFTProviderAggregator.SupplierPublic}>
    
    prepare(
        supplier: AuthAccount,
    ) {
        // Claim the aggregated NFT provider capability published by the manager
        let supplierFactoryCapability = supplier.inbox.claim<
            &NFTProviderAggregator.Aggregator{NFTProviderAggregator.SupplierFactory}>(
            capabilityPublicationID,
            provider: manager
            ) ?? panic("Could not claim capability!")

        // Borrow a reference from the capability
        self.supplierFactoryRef = supplierFactoryCapability.borrow()
            ?? panic("Could not borrow capability!")

        // Create Supplier resource and save to storage
        supplier.save(
            <-self.supplierFactoryRef.createSupplier(),
            to: NFTProviderAggregator.SupplierStoragePath
            )

        // Create supplier public capability
        self.supplierPublicCapability = supplier.link<&NFTProviderAggregator.Supplier{NFTProviderAggregator.SupplierPublic}>(
            NFTProviderAggregator.SupplierPublicPath,
            target: NFTProviderAggregator.SupplierStoragePath
        ) ?? panic("Could not link Supplier capability!")
    }
    
    post {
        // Verify that the supplier now owns a Supplier resource
        self.supplierPublicCapability.borrow()?.getAggregatorUUID() == self.supplierFactoryRef.uuid: "Supplier resource was not created!"
    }
}
