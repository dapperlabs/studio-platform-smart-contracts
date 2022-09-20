import PDS from "../../contracts/PDS.cdc"
import {{.CollectibleNFTName}} from 0x{{.CollectibleNFTAddress}}

transaction (distId: UInt64, nftIDs: [UInt64]) {
    prepare(pds: AuthAccount) {
        let cap = pds.borrow<&PDS.DistributionManager>(from: PDS.DistManagerStoragePath) ?? panic("pds does not have Dist manager")
        cap.withdraw(distId: distId, nftIDs: nftIDs, escrowCollectionPublic: {{.CollectibleNFTName}}.CollectionPublicPath)
    }
}
