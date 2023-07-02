import Golazos from "../../contracts/Golazos.cdc"
import MetadataViews from 0x{{.MetadataViewsAddress}}


pub fun main(address: Address, id: UInt64): MetadataViews.Traits {
    let account = getAccount(address)

    let collectionRef = account.getCapability(Golazos.CollectionPublicPath)
                            .borrow<&{Golazos.MomentNFTCollectionPublic}>()!

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the metadata information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Traits>())!

    return view as! MetadataViews.Traits
}

