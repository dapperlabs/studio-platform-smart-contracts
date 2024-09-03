import Golazos from "Golazos"
import MetadataViews from "MetadataViews"


access(all) fun main(address: Address, id: UInt64): MetadataViews.Editions {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
        ?? panic("Could not borrow a reference of the public collection")

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the basic display information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Editions>())!

    return view as! MetadataViews.Editions
}

