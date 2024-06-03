import NonFungibleToken from "NonFungibleToken"
import DapperSport from "DapperSport"

transaction(NFTProviderPath: PrivatePath) {

    prepare(signer: auth(Capabilities) &Account) {
        if signer.getCapability<&{NonFungibleToken.Provider}>(NFTProviderPath).check() {
            return
        }
        let cap = signer.capabilities.storage.issue<&{NonFungibleToken.Provider}>(target: DapperSport.CollectionStoragePath)

        // This needs to be used to allow for PDS to withdraw
        signer.capabilities.publish(cap, at: NFTProviderPath)
    }

}
