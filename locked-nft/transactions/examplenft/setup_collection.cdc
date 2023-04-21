import NonFungibleToken from 0xf8d6e0586b0a20c7
import ExampleNFT from 0xf8d6e0586b0a20c7

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