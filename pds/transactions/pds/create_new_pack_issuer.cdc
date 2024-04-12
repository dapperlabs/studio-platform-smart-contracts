import PDS from "PDS"

transaction() {
    prepare (issuer: auth(Storage, Capabilities) &Account) {
        // Check if account already have a PackIssuer resource, if so destroy it
        if issuer.storage.borrow<&PDS.PackIssuer>(from: PDS.PackIssuerStoragePath) != nil {
            issuer.capabilities.unpublish(PDS.PackIssuerCapRecv)
            let p <- issuer.storage.load<@PDS.PackIssuer>(from: PDS.PackIssuerStoragePath)
            destroy p
        }

        issuer.storage.save(<- PDS.createPackIssuer(), to: PDS.PackIssuerStoragePath);

        issuer.capabilities.publish(
            issuer.capabilities.storage.issue<&PDS.PackIssuer>(PDS.PackIssuerStoragePath),
            at: PDS.PackIssuerCapRecv
        )
    }
}
