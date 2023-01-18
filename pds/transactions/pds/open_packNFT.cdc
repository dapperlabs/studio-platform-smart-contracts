import PDS from "../../contracts/PDS.cdc"
import {{.CollectibleNFTName}} from 0x{{.CollectibleNFTAddress}}
import NonFungibleToken from "../../contracts/NonFungibleToken.cdc"

transaction (distId: UInt64, packId: UInt64, nftContractAddrs: [Address], nftContractName: [String], nftIds: [UInt64], owner: Address, NFTProviderPath: PrivatePath) {
    prepare(pds: AuthAccount) {
        let cap = pds.borrow<&PDS.DistributionManager>(from: PDS.DistManagerStoragePath) ?? panic("pds does not have Dist manager")
        let recvAcct = getAccount(owner)
        let recv = recvAcct.getCapability({{.CollectibleNFTName}}.CollectionPublicPath).borrow<&{NonFungibleToken.CollectionPublic}>()
            ?? panic("Unable to borrow Collection Public reference for recipient")
        cap.openPackNFT(
            distId: distId,
            packId: packId,
            nftContractAddrs: nftContractAddrs,
            nftContractName: nftContractName,
            nftIds: nftIds,
            recvCap: recv,
            collectionProviderPath: NFTProviderPath,
        )
    }
}