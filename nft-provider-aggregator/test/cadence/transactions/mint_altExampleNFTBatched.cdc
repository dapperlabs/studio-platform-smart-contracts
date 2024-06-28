import NonFungibleToken from "./NonFungibleToken.cdc"
import AltExampleNFT from "./../test/cadence/contracts/AltExampleNFT.cdc"

transaction(recipient: Address, batchSize: Int) {

    let minter: &AltExampleNFT.NFTMinter

    prepare(signer: AuthAccount) {
        self.minter = signer
            .borrow<&AltExampleNFT.NFTMinter>(from: AltExampleNFT.MinterStoragePath)!
    }

    execute {
        let receiver = getAccount(recipient)
            .getCapability(AltExampleNFT.CollectionPublicPath)!
            .borrow<&{NonFungibleToken.CollectionPublic}>()!

        var i = 0
        while i < batchSize {
            self.minter.mintNFT(recipient: receiver)
            i = i + 1
        }
    }
}
