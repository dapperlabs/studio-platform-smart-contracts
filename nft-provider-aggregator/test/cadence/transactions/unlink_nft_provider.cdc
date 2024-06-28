import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import NFTProviderAggregator from "NFTProviderAggregator"

transaction() {
    prepare(
        signer: auth(LoadValue) &Account,
    ) {
        // Signer unlinks their NFT provider capability
        let exampleNFTWithdrawCapPath: StoragePath = /storage/exampleNFTProvider
        signer.storage.load<Capability<auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection>>(from: exampleNFTWithdrawCapPath)
    }
}
