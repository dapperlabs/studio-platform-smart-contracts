import NFTProviderAggregator from "NFTProviderAggregator"

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

    let supplierFactoryRef: auth(NFTProviderAggregator.Operate) &{NFTProviderAggregator.SupplierFactory}
    let supplierPublicCapability: Capability<&{NFTProviderAggregator.SupplierPublic}>

    prepare(
        supplier: auth(BorrowValue, SaveValue, ClaimInboxCapability, IssueStorageCapabilityController, PublishCapability) &Account,
    ) {
        // Claim the aggregated NFT provider capability published by the manager
        let supplierFactoryCapability = supplier.inbox.claim<
            auth(NFTProviderAggregator.Operate) &{NFTProviderAggregator.SupplierFactory}>(
            capabilityPublicationID,
            provider: manager
            ) ?? panic("Could not claim capability!")

        // Borrow a reference from the capability
        self.supplierFactoryRef = supplierFactoryCapability.borrow()
            ?? panic("Could not borrow capability!")

        // Create Supplier resource and save to storage
        supplier.storage.save(
            <-self.supplierFactoryRef.createSupplier(),
            to: NFTProviderAggregator.SupplierStoragePath
            )

        // Create supplier public capability
        self.supplierPublicCapability = supplier.capabilities.storage.issue<&{NFTProviderAggregator.SupplierPublic}>(
            NFTProviderAggregator.SupplierStoragePath
        )

        // Publish the supplier public capability
        supplier.capabilities.publish(
            self.supplierPublicCapability,
            at: NFTProviderAggregator.SupplierPublicPath,
        )
    }

    post {
        // Verify that the supplier now owns a Supplier resource
        self.supplierPublicCapability.borrow()?.getAggregatorUUID() == self.supplierFactoryRef.uuid: "Supplier resource was not created!"
    }
}
