import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import ExampleNFT from "../../contracts/ExampleNFT.cdc"

transaction(NFTProviderPath: PrivatePath) {

    prepare(signer: AuthAccount) {
        if signer.getCapability<&{NonFungibleToken.Provider}>(NFTProviderPath).check() {
            return
        }
        // This needs to be used to allow for PDS to withdraw
        signer.link<&{NonFungibleToken.Provider}>( NFTProviderPath, target: ExampleNFT.CollectionStoragePath)
    }

}
