import NonFungibleToken from "NonFungibleToken"
import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a manager account to publish a private capability to its aggregated NFT provider
/// to be later claimed by the recipient (who is trusted for the ability to use the capability itself but also for
/// potentially copying and saving it somewhere else).
///
/// @param recipient: The third-party recipient address to publish the aggregated NFT provider capability to.
/// @param capabilityPublicationID: The publication identifier of the aggregated NFT provider capability.
///
transaction(
    recipient: Address,
    capabilityPublicationID: String
    ) {

    let aggregatedNFTWithdrawCap: Capability<
        auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>

    prepare(
        manager: auth(IssueStorageCapabilityController, CopyValue, SaveValue, PublishInboxCapability) &Account,
    ) {
        // Retrieve or create aggregated NFT provider capability
        let aggregatedNFTWithdrawCapStoragePath = NFTProviderAggregator.convertPrivateToStoragePath(NFTProviderAggregator.AggregatedProviderPrivatePath)
        if let retrievedCap = manager.storage.copy<Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>>(
                from: aggregatedNFTWithdrawCapStoragePath) {
            self.aggregatedNFTWithdrawCap = retrievedCap
        } else {
            self.aggregatedNFTWithdrawCap = manager.capabilities.storage.issue<
                auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>(
                NFTProviderAggregator.AggregatorStoragePath)
            self.aggregatedNFTWithdrawCap.check()

            manager.storage.save(self.aggregatedNFTWithdrawCap, to: aggregatedNFTWithdrawCapStoragePath)
        }
        assert(self.aggregatedNFTWithdrawCap.check(), message: "Capability is invalid")


        // Publish the aggregated NFT provider capability to recipient
        manager.inbox.publish(
            self.aggregatedNFTWithdrawCap,
            name: capabilityPublicationID,
            recipient: recipient)
    }
}
