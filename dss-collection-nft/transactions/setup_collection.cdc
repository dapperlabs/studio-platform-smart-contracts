import DSSCollection from 0xf8d6e0586b0a20c7
import NonFungibleToken from 0xf8d6e0586b0a20c7

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&DSSCollection.Collection>(from: DSSCollection.CollectionStoragePath) != nil {
            return;
        }
        let collection <- DSSCollection.createEmptyCollection();
        signer.save(<-collection, to: DSSCollection.CollectionStoragePath);
        signer.link<&DSSCollection.Collection{NonFungibleToken.CollectionPublic, DSSCollection.DSSCollectionNFTCollectionPublic}>(
            DSSCollection.CollectionPublicPath,
            target: DSSCollection.CollectionStoragePath
        ) ?? panic("Could not link DSSCollection.Collection Pub Path");
    }
}