import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import NFTProviderAggregator from "NFTProviderAggregator"

// Delete all withdraw capabilities targeting the default ExampleNFT collection storage path
transaction() {
    var countDeleted: Int
    prepare(
        signer: auth(LoadValue, Capabilities) &Account,
    ) {
        self.countDeleted  = 0
        for controller in signer.capabilities.storage.getControllers(forPath: ExampleNFT.CollectionStoragePath) {
            if controller.tag == "nft-provider-aggregator" {
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
