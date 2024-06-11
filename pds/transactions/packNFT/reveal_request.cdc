import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"
import NonFungibleToken from "NonFungibleToken"


transaction(revealID: UInt64, openRequest: Bool) {
    prepare(owner: auth(BorrowValue) &Account) {
        let collectionRef = owner.storage.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
        let packNFT = collectionRef.borrowNFT(id: revealID) as! auth(NonFungibleToken.Update) &{IPackNFT.NFT}!
        packNFT.reveal(openRequest: openRequest)
    }
}
