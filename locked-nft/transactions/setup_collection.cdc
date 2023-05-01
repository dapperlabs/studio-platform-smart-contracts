import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import NFTLocker from "../contracts/NFTLocker.cdc"

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&NFTLocker.Collection>(from: NFTLocker.CollectionStoragePath) == nil {

            let collection <- NFTLocker.createEmptyCollection()
            signer.save(<-collection, to: NFTLocker.CollectionStoragePath)
            signer.link<&NFTLocker.Collection{NFTLocker.LockedCollection}>(
                NFTLocker.CollectionPublicPath,
                target: NFTLocker.CollectionStoragePath
            )
        }
    }
}