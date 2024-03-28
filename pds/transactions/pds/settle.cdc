import PDS from "PDS"
import ExampleNFT from "ExampleNFT"

transaction (distId: UInt64, nftIDs: [UInt64]) {
    prepare(pds: auth(BorrowValue) &Account) {
        let cap = pds.storage.borrow<&PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")
        cap.withdraw(
            distId: distId,
            nftIDs: nftIDs,
            escrowCollectionPublic: PublicPath(identifier: "cadenceExampleNFTCollection")!,
        )
    }
}
