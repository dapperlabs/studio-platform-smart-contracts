import NonFungibleToken from 0xf8d6e0586b0a20c7
import LockedNFT from 0xf8d6e0586b0a20c7

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&LockedNFT.Collection>(from: LockedNFT.CollectionStoragePath) == nil {

            let collection <- LockedNFT.createEmptyCollection()
            signer.save(<-collection, to: LockedNFT.CollectionStoragePath)
            signer.link<&LockedNFT.Collection{LockedNFT.LockedCollection}>(
                LockedNFT.CollectionPublicPath,
                target: LockedNFT.CollectionStoragePath
            )
        }
    }
}