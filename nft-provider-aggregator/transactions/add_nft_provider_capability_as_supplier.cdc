import NonFungibleToken from "NonFungibleToken"
import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a supplier to add a NFT provider capability to the parent Aggregator resource.
///
/// @param nftProviderPrivatePathID: The private path ID of the NFT provider to add - the ID is the path without
/// the domain prefix (e.g., for "/private/exampleNFTProvider", the ID is "exampleNFTProvider").
/// @param nftCollectionStoragePathID: The storage path ID of the NFT collection to add - the ID is the path without
/// the domain prefix (e.g., for "/storage/exampleNFTCollection", the ID is "exampleNFTCollection").
///
transaction(
    nftWithdrawCapStoragePathID: String,
    nftCollectionStoragePathID: String
    ) {

    let nftWithdrawCapStoragePath: StoragePath
    let nftCollectionStoragePath: StoragePath
    let supplierRef: auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Supplier
    let nftWithdrawCapability: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>

    prepare(
        supplier: auth(BorrowValue, Storage, Capabilities) &Account,
    ) {
        // Convert provided string paths
        self.nftWithdrawCapStoragePath = StoragePath(identifier: nftWithdrawCapStoragePathID)
            ?? panic("Provided NFT provider private path has invalid format!")
        self.nftCollectionStoragePath = StoragePath(identifier: nftCollectionStoragePathID)
            ?? panic("Provided NFT collection storage path has invalid format!")

        // Retrieve or create NFT provider capability
        if let retrievedCap = supplier.storage.copy<Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>>(
                from: self.nftWithdrawCapStoragePath) {
            self.nftWithdrawCapability = retrievedCap
        } else {
            self.nftWithdrawCapability = supplier.capabilities.storage.issue<
                auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
                self.nftCollectionStoragePath)
            supplier.capabilities.storage.getController(byCapabilityID: self.nftWithdrawCapability.id)!.setTag("nft-provider-aggregator")
            supplier.storage.save(self.nftWithdrawCapability, to: self.nftWithdrawCapStoragePath)
        }

        // Create reference to Supplier resource from storage
        self.supplierRef = supplier.storage.borrow<auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath)!
    }

    execute {
        // Add NFT provider capability
        let collectionUUID = self.supplierRef.addNFTWithdrawCapability(self.nftWithdrawCapability)
        log(collectionUUID)
    }
}
