import PDS from "PDS"

transaction (distId: UInt64, title: String, metadata: {String: String}) {
    prepare(issuer: auth(BorrowValue) &Account) {
        // get pack issuer reference
        let issuerRef = issuer.storage.borrow<auth(PDS.Operate) &PDS.PackIssuer>(from: PDS.PackIssuerStoragePath)
            ?? panic ("issuer does not have PackIssuer resource")

        // update distribution
        issuerRef.updateDist(
            distId: distId,
            title: title,
            metadata: metadata,
        )
    }
}
