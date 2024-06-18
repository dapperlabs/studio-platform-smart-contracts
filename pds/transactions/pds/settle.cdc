import PDS from "PDS"
import ExampleNFT from "ExampleNFT"

transaction (NFTProviderPath: String, distId: UInt64, nftIDs: [UInt64]) {
    prepare(pds: auth(BorrowValue) &Account) {
        let cap = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        cap.withdraw(
            distId: distId,
            nftIDs: nftIDs,
            escrowCollectionPublic: PublicPath(identifier: NFTProviderPath)!,
        )
    }
}
