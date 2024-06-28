import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a manager to add a NFT provider capability to the Aggregator resource.
///
/// @param nftProviderPrivatePathID: The private path ID of the NFT provider to add - the ID is the path without
/// the domain prefix (e.g., for "/private/exampleNFTProvider", the ID is "exampleNFTProvider").
/// @param nftCollectionStoragePathID: The storage path IDof the NFT collection to add - the ID is the path without
/// the domain prefix (e.g., for "/storage/exampleNFTCollection", the ID is "exampleNFTCollection").
///
transaction(
    nftProviderPrivatePathID: String,
    nftCollectionStoragePathID: String
    ) {

    let privatePath: PrivatePath
    let storagePath: StoragePath
    let aggregatorRef: &NFTProviderAggregator.Aggregator
    let nftProviderCapability: Capability<
        &AnyResource{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>

    prepare(
        manager: AuthAccount,
    ) {
        // Convert provided string paths
        self.privatePath = PrivatePath(identifier: nftProviderPrivatePathID)
            ?? panic("Provided NFT provider private path has invalid format!")
        self.storagePath = StoragePath(identifier: nftCollectionStoragePathID)
            ?? panic("Provided NFT collection storage path has invalid format!")

        // Retrieve or create NFT provider capability
        let retrievedCap = manager.getCapability<
            &AnyResource{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(
                self.privatePath)
        if retrievedCap.check() {
            self.nftProviderCapability = retrievedCap
        }
        else {
            self.nftProviderCapability = manager.link<
                &AnyResource{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>(
                self.privatePath,
                target: self.storagePath)!
        }

        // Create reference to Aggregator resource from storage
        self.aggregatorRef = manager.borrow<&NFTProviderAggregator.Aggregator>(
            from: NFTProviderAggregator.AggregatorStoragePath)!
    }

    execute {
        // Add NFT provider capability
        self.aggregatorRef.addNFTProviderCapability(nftProviderCapability:  self.nftProviderCapability)
    }
}
