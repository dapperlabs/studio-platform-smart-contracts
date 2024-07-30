import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import NFTProviderAggregator from "NFTProviderAggregator"

/// Delete all withdraw capabilities targeting the default ExampleNFT collection storage path
///
/// @param withdrawCapabilityTag: The tag to set on the capability controller to keep track of the capability being
/// supplied to a NFT provider aggregator and faciliate revokation when needed
///
transaction(
        withdrawCapabilityTag: String
) {
    var countDeleted: Int
    prepare(
        signer: auth(LoadValue, Capabilities) &Account,
    ) {
        self.countDeleted  = 0
        for controller in signer.capabilities.storage.getControllers(forPath: ExampleNFT.CollectionStoragePath) {
            if controller.tag == withdrawCapabilityTag {
                assert(controller.borrowType == Type<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(),
                    message: "Unexpected capability type")

                controller.delete()
                self.countDeleted = self.countDeleted + 1
            }
        }
    }

    post {
        self.countDeleted == 1: "expected to delete one capability, got ".concat(self.countDeleted.toString())
    }
}
