import PDS from "PDS"
import PackNFT from "PackNFT"
import ExampleNFT from "ExampleNFT"
import NonFungibleToken from "NonFungibleToken"

transaction (
    distId: UInt64,
    packId: UInt64,
    nftContractAddrs: [Address],
    nftContractNames: [String],
    nftIds: [UInt64],
    salt: String,
    owner: Address,
    openRequest: Bool,
    collectionStoragePath: StoragePath
) {
    prepare(pds: auth(BorrowValue) &Account) {
        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        let cap = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        let p = PackNFT.borrowPackRepresentation(id: packId)
            ?? panic ("No such pack")

        if openRequest && p.status == PackNFT.Status.Revealed {
            let recvAcct = getAccount(owner)
            let recv = recvAcct.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(collectionData.publicPath)
                ?? panic("Unable to borrow Collection Public reference for recipient")
            cap.openPackNFT(
                distId: distId,
                packId: packId,
                nftContractAddrs: nftContractAddrs,
                nftContractNames: nftContractNames,
                nftIds: nftIds,
                recvCap: recv,
                collectionStoragePath: collectionStoragePath
            )
        } else {
            cap.revealPackNFT(
                distId: distId,
                packId: packId,
                nftContractAddrs: nftContractAddrs,
                nftContractNames: nftContractNames,
                nftIds: nftIds,
                salt: salt
            )
        }
    }
}
