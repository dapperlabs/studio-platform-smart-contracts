import NonFungibleToken from "../../../contracts/NonFungibleToken.cdc"
import ExampleNFT from "../../../contracts/ExampleNFT.cdc"
import NFTProviderAggregator from "../../../contracts/NFTProviderAggregator.cdc"

transaction() {
    
    prepare(
        signer: AuthAccount,
    ) {
        // Signer unlinks their NFT provider capability
        let exampleNFTProviderPath: PrivatePath = /private/exampleNFTProvider
        signer.unlink(exampleNFTProviderPath)
    }
}
