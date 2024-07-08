import PDS from "PDS"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

transaction (distId: UInt64, nftIDs: [UInt64]) {
    prepare(pds: auth(BorrowValue) &Account) {
        let collectionData = ExampleNFT.resolveContractView(
            resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")
        let cap = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        cap.withdraw(
            distId: distId,
            nftIDs: nftIDs,
            escrowCollectionPublic: collectionData.publicPath,
        )
    }
}
