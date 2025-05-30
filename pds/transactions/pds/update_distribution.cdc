import PDS from "PDS"

transaction (distId: UInt64, title: String, metadata: {String: String}) {
    prepare(pds: auth(BorrowValue) &Account) {
        // get pack issuer reference
        let managerRef = pds.storage.borrow<auth(PDS.Operate) &PDS.DistributionManager>(from: PDS.DistManagerStoragePath)
            ?? panic("pds does not have Dist manager")

        // update distribution
        managerRef.updateDist(
            distId: distId,
            title: title,
            metadata: metadata,
        )
    }
}
