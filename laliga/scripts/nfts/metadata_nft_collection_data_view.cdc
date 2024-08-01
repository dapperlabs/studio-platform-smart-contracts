import Golazos from "Golazos"
import MetadataViews from "MetadataViews"

 access(all) struct NFTCollectionData {
        access(all) let storagePath: StoragePath

        access(all) let publicPath: PublicPath

        access(all) let publicCollection: Type

        access(all) let publicLinkedType: Type

        init(
            storagePath: StoragePath,
            publicPath: PublicPath,
            publicCollection: Type,
            publicLinkedType: Type,
        ) {
            self.storagePath=storagePath
            self.publicPath=publicPath
            self.publicCollection=publicCollection
            self.publicLinkedType=publicLinkedType
        }
    }

access(all) fun main(address: Address, id: UInt64): NFTCollectionData {
    let account = getAccount(address)

    let collectionRef = account.capabilities.borrow<&Golazos.Collection>(Golazos.CollectionPublicPath)
                            ?? panic("Could not borrow a reference of the public collection")

    let nft = collectionRef.borrowMomentNFT(id: id)!
    
    // Get the NFTCollectionData information for this NFT
    let data = nft.resolveView(Type<MetadataViews.NFTCollectionData>())! as! MetadataViews.NFTCollectionData
    return NFTCollectionData(storagePath: data.storagePath,
            publicPath: data.publicPath,
            publicCollection: data.publicCollection,
            publicLinkedType: data.publicLinkedType)
}

