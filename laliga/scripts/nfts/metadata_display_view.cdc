import Golazos from "Golazos"
import MetadataViews from "MetadataViews"

access(all) struct NFT {
    access(all) let name: String
    access(all) let description: String
    access(all) let thumbnail: String

    init(
        name: String,
        description: String,
        thumbnail: String,
    ) {
        self.name = name
        self.description = description
        self.thumbnail = thumbnail
    }
}

access(all) fun main(address: Address, id: UInt64): NFT {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<
        &Golazos.Collection>(Golazos.CollectionPublicPath)
        ?? panic("Could not borrow a reference of the public collection")

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the basic display information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Display>())!

    let display = view as! MetadataViews.Display

    return NFT(
        name: display.name,
        description: display.description,
        thumbnail: display.thumbnail.uri()
    )
}

