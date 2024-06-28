import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a supplier to add a NFT provider capability to the parent Aggregator resource.
///
/// @param nftProviderPrivatePathID: The private path ID of the NFT provider to add - the ID is the path without
/// the domain prefix (e.g., for "/private/exampleNFTProvider", the ID is "exampleNFTProvider").
/// @param nftCollectionStoragePathID: The storage path ID of the NFT collection to add - the ID is the path without
/// the domain prefix (e.g., for "/storage/exampleNFTCollection", the ID is "exampleNFTCollection").
///
transaction(
    nftProviderPrivatePathID: String,
    nftCollectionStoragePathID: String
    ) {
    
    let privatePath: PrivatePath
    let storagePath: StoragePath
    let supplierRef: &NFTProviderAggregator.Supplier
    let nftProviderCapability: Capability<
        &AnyResource{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>

    prepare(
        supplier: AuthAccount,
    ) {
        // Convert provided string paths
        self.privatePath = PrivatePath(identifier: nftProviderPrivatePathID)
            ?? panic("Provided NFT provider private path has invalid format!")
        self.storagePath = StoragePath(identifier: nftCollectionStoragePathID)
            ?? panic("Provided NFT collection storage path has invalid format!")

        // Retrieve or create NFT provider capability
        let retrievedCap = supplier.getCapability<
            &AnyResource{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(
                self.privatePath)
        if retrievedCap.check() {
            self.nftProviderCapability = retrievedCap
        }
        else {
            self.nftProviderCapability = supplier.link<
                &AnyResource{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(
                self.privatePath,
                target: self.storagePath)!
        }

        // Create reference to Supplier resource from storage
        self.supplierRef = supplier.borrow<&NFTProviderAggregator.Supplier>(
            from: NFTProviderAggregator.SupplierStoragePath)!
    }

    execute {
        // Add NFT provider capability
        self.supplierRef.addNFTProviderCapability(nftProviderCapability: self.nftProviderCapability)
    }
}
