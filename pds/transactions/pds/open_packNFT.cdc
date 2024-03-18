import PDS from "PDS"
import ExampleNFT from "ExampleNFT"
import NonFungibleToken from "NonFungibleToken"

transaction (
    distId: UInt64,
    packId: UInt64,
    nftContractAddrs: [Address],
    nftContractName: [String],
    nftIds: [UInt64],
    owner: Address,
    collectionStoragePath: StoragePath
) {
    prepare(pds: auth(BorrowValue) &Account) {
        let cap = pds.storage.borrow<&PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        let recvAcct = getAccount(owner)
        let recv = recvAcct.capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(PublicPath(identifier: "cadenceExampleNFTCollection")!)
            ?? panic("Unable to borrow Collection Public reference for recipient")

        cap.openPackNFT(
            distId: distId,
            packId: packId,
            nftContractAddrs: nftContractAddrs,
            nftContractName: nftContractName,
            nftIds: nftIds,
            recvCap: recv,
            collectionStoragePath: collectionStoragePath,
        )
    }
}
