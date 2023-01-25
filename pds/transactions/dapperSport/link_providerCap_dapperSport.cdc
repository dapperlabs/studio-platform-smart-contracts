import NonFungibleToken from 0x{{.NonFungibleToken}}
import {{.DapperSportContract}} from 0x{{.DapperSportAddress}}

transaction(NFTProviderPath: PrivatePath) {

    prepare(signer: AuthAccount) {
        if signer.getCapability<&{NonFungibleToken.Provider}>(NFTProviderPath).check() {
            return
        }
        // This needs to be used to allow for PDS to withdraw
        signer.link<&{NonFungibleToken.Provider}>( NFTProviderPath, target: {{.DapperSportContract}}.CollectionStoragePath)
    }

}
