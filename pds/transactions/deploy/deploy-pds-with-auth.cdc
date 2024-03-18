// This transactions deploys a contract with init args
//
transaction(
    contractName: String,
    code: String,
    PackIssuerStoragePath: StoragePath,
    PackIssuerCapRecv: PublicPath,
    DistCreatorStoragePath: StoragePath,
    DistManagerStoragePath: StoragePath,
    version: String,
) {
    prepare(owner: auth(AddContract) &Account) {
        let existingContract = owner.contracts.get(name: contractName)

        if (existingContract == nil) {
            log("no contract")
            owner.contracts.add(
                name: contractName,
                code: code.decodeHex(),
                PackIssuerStoragePath: PackIssuerStoragePath,
                PackIssuerCapRecv: PackIssuerCapRecv,
                DistCreatorStoragePath: DistCreatorStoragePath,
                DistManagerStoragePath: DistManagerStoragePath,
                version: version,
            )
        } else {
            log("has contract")
            owner.contracts.add(name: contractName, code: code.decodeHex())
        }
    }
}
