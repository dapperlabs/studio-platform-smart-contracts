import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

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

    let aggregatedNFTProviderCap: Capability<
        &AnyResource{NonFungibleToken.Provider}>

    prepare(
        manager: AuthAccount,
    ) {
        // Retrieve or create aggregated NFT provider capability
        let retrievedCap = manager.getCapability<
            &AnyResource{NonFungibleToken.Provider}>(
            NFTProviderAggregator.AggregatedProviderPrivatePath)
        if retrievedCap.check(){
            self.aggregatedNFTProviderCap = retrievedCap
        }
        else {
            self.aggregatedNFTProviderCap = manager.link<
                &AnyResource{NonFungibleToken.Provider}>(
                NFTProviderAggregator.AggregatedProviderPrivatePath,
                target: NFTProviderAggregator.AggregatorStoragePath)!
        }

        // Publish the aggregated NFT provider capability to recipient
        manager.inbox.publish(
            self.aggregatedNFTProviderCap,
            name: capabilityPublicationID,
            recipient: recipient)
    }
}
 