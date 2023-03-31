import EnglishPremierLeague from "./EnglishPremierLeague.cdc"
import MetadataViews from 0xMETADATAVIEWSADDRESS

pub struct NFT {
    pub let name: String
    pub let description: String
    pub let thumbnail: String

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

pub fun main(address: Address, id: UInt64): NFT {
    let account = getAccount(address)

    let collectionRef = account.getCapability(EnglishPremierLeague.CollectionPublicPath)
                            .borrow<&{EnglishPremierLeague.MomentNFTCollectionPublic}>()!

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

