import Golazos from "../../contracts/Golazos.cdc"
import MetadataViews from 0x{{.MetadataViewsAddress}}


pub fun main(address: Address, id: UInt64): MetadataViews.Editions {
    let account = getAccount(address)

    let collectionRef = account.getCapability(Golazos.CollectionPublicPath)
                            .borrow<&{Golazos.MomentNFTCollectionPublic}>()!

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the basic display information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Editions>())!

    return view as! MetadataViews.Editions
}

