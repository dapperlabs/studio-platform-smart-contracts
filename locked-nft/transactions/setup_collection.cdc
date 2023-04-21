import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import LockedNFT from "../contracts/LockedNFT.cdc"

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