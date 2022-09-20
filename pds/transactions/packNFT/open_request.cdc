import PackNFT from "../../contracts/PackNFT.cdc"
import IPackNFT from "../../contracts/IPackNFT.cdc"

transaction(revealID: UInt64) {
    prepare(owner: AuthAccount) {
        let collectionRef = owner.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
        collectionRef.borrowPackNFT(id: revealID)!.open()
    }
}
