import PackNFT from "../../contracts/PackNFT.cdc"
import IPackNFT from "../../contracts/IPackNFT.cdc"

transaction(revealID: UInt64, openRequest: Bool) {
    prepare(owner: AuthAccount) {
        let collectionRef = owner.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
        collectionRef.borrowPackNFT(id: revealID)!.reveal(openRequest: openRequest)
    }
}
