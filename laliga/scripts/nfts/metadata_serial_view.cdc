import Golazos from "../../contracts/Golazos.cdc"
import MetadataViews from 0x{{.MetadataViewsAddress}}


pub fun main(address: Address, id: UInt64): UInt64 {
    let account = getAccount(address)

    let collectionRef = account.getCapability(Golazos.CollectionPublicPath)
                            .borrow<&{Golazos.MomentNFTCollectionPublic}>()!

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the basic display information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Serial>())!

    let serial = view as! MetadataViews.Serial

    return serial.number
}

