import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

/// Transaction signed by a manager account to unpublish a capability that was previously
/// published and that hasn't been claimed yet by the intended recipient account.
///
/// @param capabilityPublicationID: The publication identifier of capability to unpublish.
///
transaction(
    capabilityPublicationID: String
    ) {

    prepare(
        manager: AuthAccount,
    ) {
        // Unpublish capability
        manager.inbox.unpublish(capabilityPublicationID)  
    }
}
