import DapperSport from "../../contracts/DapperSport.cdc"
import MetadataViews from 0xMETADATAVIEWSADDRESS

 pub struct NFTCollectionData {
        pub let storagePath: StoragePath

        pub let publicPath: PublicPath

        pub let providerPath: PrivatePath

        pub let publicCollection: Type

        pub let publicLinkedType: Type

        pub let providerLinkedType: Type

        init(
            storagePath: StoragePath,
            publicPath: PublicPath,
            providerPath: PrivatePath,
            publicCollection: Type,
            publicLinkedType: Type,
            providerLinkedType: Type
        ) {
            self.storagePath=storagePath
            self.publicPath=publicPath
            self.providerPath = providerPath
            self.publicCollection=publicCollection
            self.publicLinkedType=publicLinkedType
            self.providerLinkedType = providerLinkedType
        }
    }

pub fun main(address: Address, id: UInt64): NFTCollectionData {
    let account = getAccount(address)

    let collectionRef = account.getCapability(DapperSport.CollectionPublicPath)
                            .borrow<&{DapperSport.MomentNFTCollectionPublic}>()!

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the NFTCollectionData information for this NFT
    let data = nft.resolveView(Type<MetadataViews.NFTCollectionData>())! as! MetadataViews.NFTCollectionData
    return NFTCollectionData(storagePath: data.storagePath,
            publicPath: data.publicPath,
            providerPath: data.providerPath,
            publicCollection: data.publicCollection,
            publicLinkedType: data.publicLinkedType,
            providerLinkedType: data.providerLinkedType)
}

