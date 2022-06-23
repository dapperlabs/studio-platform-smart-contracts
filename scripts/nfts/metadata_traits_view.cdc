import DapperSport from "../../contracts/DapperSport.cdc"
import MetadataViews from 0xMETADATAVIEWSADDRESS


pub fun main(address: Address, id: UInt64): MetadataViews.Traits {
    let account = getAccount(address)

    let collectionRef = account.getCapability(DapperSport.CollectionPublicPath)
                            .borrow<&{DapperSport.MomentNFTCollectionPublic}>()!

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the metadata information for this NFT
    let view = nft.resolveView(Type<MetadataViews.Traits>())!

    return view as! MetadataViews.Traits
}

