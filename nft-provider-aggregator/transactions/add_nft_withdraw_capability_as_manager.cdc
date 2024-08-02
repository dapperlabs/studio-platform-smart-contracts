import NonFungibleToken from "NonFungibleToken"
import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a manager to add a NFT provider capability to the Aggregator resource.
///
/// @param nftWithdrawCapStoragePathID: The storage  path ID of the NFT provider to add - the ID is the path without
/// the domain prefix (e.g., for "/storage/exampleNFTProvider", the ID is "exampleNFTProvider").
/// @param nftCollectionStoragePathID: The storage path ID of the NFT collection to add - the ID is the path without
/// the domain prefix (e.g., for "/storage/exampleNFTCollection", the ID is "exampleNFTCollection").
/// @param withdrawCapabilityTag: The tag to set on the capability controller to keep track of the capability being
/// supplied to a NFT provider aggregator and faciliate revokation when needed
///
transaction(
    nftWithdrawCapStoragePathID: String,
    nftCollectionStoragePathID: String,
    withdrawCapabilityTag: String
    ) {

    let nftWithdrawCapStoragePath: StoragePath
    let nftCollectionStoragePath: StoragePath
    let aggregatorRef: auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Aggregator
    let nftWithdrawCapability: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>

    prepare(
        manager: auth(CopyValue, BorrowValue, SaveValue, GetStorageCapabilityController, IssueStorageCapabilityController) &Account,
    ) {
        // Convert provided string paths
        self.nftWithdrawCapStoragePath = StoragePath(identifier: nftWithdrawCapStoragePathID)
            ?? panic("Provided NFT provider private path has invalid format!")
        self.nftCollectionStoragePath = StoragePath(identifier: nftCollectionStoragePathID)
            ?? panic("Provided NFT collection storage path has invalid format!")

        // Retrieve or create NFT provider capability
        if let retrievedCap = manager.storage.copy<Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider, NonFungibleToken.CollectionPublic}>>(
                from: self.nftWithdrawCapStoragePath) {
            self.nftWithdrawCapability = retrievedCap
        } else {
            self.nftWithdrawCapability = manager.capabilities.storage.issue<
                auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(
                self.nftCollectionStoragePath)
            manager.capabilities.storage.getController(byCapabilityID: self.nftWithdrawCapability.id)!.setTag(withdrawCapabilityTag)
            manager.storage.save(self.nftWithdrawCapability, to: self.nftWithdrawCapStoragePath)
        }

        // Create reference to Aggregator resource from storage
        self.aggregatorRef = manager.storage.borrow<auth(NFTProviderAggregator.Operate) &NFTProviderAggregator.Aggregator>(
            from: NFTProviderAggregator.AggregatorStoragePath)!
    }

    execute {
        // Add NFT provider capability
        let collectionUUID = self.aggregatorRef.addNFTWithdrawCapability(self.nftWithdrawCapability)
        log(collectionUUID)
    }
}
