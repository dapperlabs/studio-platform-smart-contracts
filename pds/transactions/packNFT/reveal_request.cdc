import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"

transaction(revealID: UInt64, openRequest: Bool) {
    prepare(owner: auth(BorrowValue) &Account) {
        let collectionRef = owner.storage.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
        collectionRef.borrowPackNFT(id: revealID)!.reveal(openRequest: openRequest)
    }
}
