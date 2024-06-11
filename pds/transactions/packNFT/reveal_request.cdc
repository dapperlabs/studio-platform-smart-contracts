import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"
import NonFungibleToken from "NonFungibleToken"


transaction(revealID: UInt64, openRequest: Bool) {
    prepare(owner: auth(BorrowValue) &Account) {
        let collectionRef = owner.storage.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath)!
        let packNFT = collectionRef.borrowNFT(revealID) as! auth(NonFungibleToken.Update) &{PackNFT.NFT}!
        packNFT.reveal(openRequest: openRequest)
    }
}
