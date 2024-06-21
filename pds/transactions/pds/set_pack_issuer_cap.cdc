import PDS from "PDS"

transaction (issuer: Address) {
    prepare(pds: auth(Capabilities) &Account) {
        let cap = pds.capabilities.storage.issue<&PDS.DistributionCreator>(PDS.DistCreatorStoragePath)
        if !cap.check() {
            panic("cannot borrow such capability")
        }

        let setCapRef = getAccount(issuer).capabilities.borrow<&PDS.PackIssuer>(PDS.PackIssuerCapRecv)
            ?? panic("no cap for setting distCap")
        setCapRef.setDistCap(cap: cap);
    }
}
