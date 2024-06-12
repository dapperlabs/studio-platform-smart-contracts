import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"
import NonFungibleToken from "NonFungibleToken"

transaction(revealID: UInt64, openRequest: Bool) {
    prepare(owner: auth(BorrowValue) &Account) {
        let collectionRef = owner.storage.borrow<auth(NonFungibleToken.Update) &PackNFT.Collection>(from: PackNFT.CollectionStoragePath)
            ?? panic("could not borrow authorized collection")
        collectionRef.reveal(id: revealID, openRequest: openRequest)
    }
}
