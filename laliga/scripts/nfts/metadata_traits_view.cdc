import Golazos from "Golazos"
import MetadataViews from "MetadataViews"


access(all) fun main(address: Address, id: UInt64): MetadataViews.Traits {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
        ?? panic("Could not borrow a reference of the public collection")
    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the metadata information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Traits>())!

    return view as! MetadataViews.Traits
}

