import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"
import PackNFT from "../../contracts/PackNFT.cdc"

transaction() {
    prepare (issuer: AuthAccount) {
        
        // Check if account already have a PackIssuer resource, if so destroy it
        if issuer.borrow<&PackNFT.Collection>(from: PackNFT.CollectionStoragePath) != nil {
            issuer.unlink(PackNFT.CollectionPublicPath)
            let p <- issuer.load<@PackNFT.Collection>(from: PackNFT.CollectionStoragePath) 
            destroy p
        }
        
        issuer.save(<- PackNFT.createEmptyCollection(), to: PackNFT.CollectionStoragePath);
        
        issuer.link<&{NonFungibleToken.CollectionPublic}>(PackNFT.CollectionPublicPath, target: PackNFT.CollectionStoragePath)
        ??  panic("Could not link Collection Pub Path");
    } 
}
 
