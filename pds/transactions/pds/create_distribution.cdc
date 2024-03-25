import PDS from "PDS"
import {{.PackNFTName}} from "PackNFT"
import IPackNFT from "IPackNFT"
import NonFungibleToken from "NonFungibleToken"

transaction(title: String, metadata: {String: String}) {
    prepare (issuer: auth(BorrowValue, Capabilities) &Account) {

        let i = issuer.storage.borrow<auth(PDS.DistCreation) &PDS.PackIssuer>(from: PDS.PackIssuerStoragePath)
            ?? panic ("issuer does not have PackIssuer resource")

        // issuer must have a PackNFT collection
        let withdrawCap = issuer.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Collection}>(StoragePath(identifier: "cadenceExampleNFTCollection")!);
        let operatorCap = issuer.capabilities.storage.issue<auth(IPackNFT.Operatable) &{IPackNFT.IOperator}>({{.PackNFTName}}.OperatorStoragePath);
        assert(withdrawCap.check(), message:  "cannot borrow withdraw capability")
        assert(operatorCap.check(), message:  "cannot borrow operator capability")

        let sc <- PDS.createSharedCapabilities ( withdrawCap: withdrawCap, operatorCap: operatorCap )
        i.createDist(sharedCap: <-sc, title: title, metadata: metadata)
    }
}
