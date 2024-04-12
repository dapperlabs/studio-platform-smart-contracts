import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"

transaction(revealID: UInt64) {
    prepare(owner: AuthAccount) {
        let collectionRef = owner.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
        collectionRef.borrowPackNFT(id: revealID)!.open()
    }
}
