import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by the intended third-party recipient account of an aggregated NFT
/// provider capability previously published by a manager.
/// Note: A published capability can only be claimed once.
///
/// @param manager: The address of the manager that published the aggregated NFT provider capability.
/// @param capabilityPublicationID: The publication identifier of the aggregated NFT provider capability.
///
transaction(
    manager: Address,
    capabilityPublicationID: String,
    ) {

    prepare(
        signer: AuthAccount,
    ) {
        // Claim the aggregated NFT provider capability published by the manager
        let capability = signer.inbox.claim<
            &AnyResource{NonFungibleToken.Provider}>(
            capabilityPublicationID,
            provider: manager
            ) ?? panic("Could not claim capability!")

        // Save capability to storage
        // Note: It is not possible to store claimed capabilities in CapabilityPath at the moment
        // (like the link() method does, so we store in StoragePath)
        signer.save(capability, to: NFTProviderAggregator.AggregatorStoragePath)
    }
}
