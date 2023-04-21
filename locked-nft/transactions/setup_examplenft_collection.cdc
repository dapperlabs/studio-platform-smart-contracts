import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import ExampleNFT from 0xEXAMPLENFTADDRESS


transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath) == nil {

            let collection <- ExampleNFT.createEmptyCollection()
            signer.save(<-collection, to: ExampleNFT.CollectionStoragePath)
            signer.link<&ExampleNFT.Collection{NonFungibleToken.CollectionPublic, ExampleNFT.ExampleNFTCollectionPublic}>(
                ExampleNFT.CollectionPublicPath,
                target: ExampleNFT.CollectionStoragePath
            )
        }
    }
}