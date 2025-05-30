import PDS from "PDS"
import PackNFT from "PackNFT"
import IPackNFT from "IPackNFT"
import NonFungibleToken from "NonFungibleToken"

transaction(nftWithdrawCapPath: StoragePath, title: String, metadata: {String: String}) {
    prepare (issuer: auth(Storage, Capabilities) &Account) {
        // get pack issuer reference
        let issuerRef = issuer.storage.borrow<auth(PDS.Operate) &PDS.PackIssuer>(from: PDS.PackIssuerStoragePath)
            ?? panic ("issuer does not have PackIssuer resource")

        // get withdraw capability from issuer
        let withdrawCap = issuer.storage.copy<
        Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>>(from: nftWithdrawCapPath)!
        assert(withdrawCap.check(), message:  "cannot get copy of withdraw capability")

        // get operator capability from issuer
        let operatorCap = issuer.capabilities.storage.issue<auth(IPackNFT.Operate) &{IPackNFT.IOperator}>(PackNFT.OperatorStoragePath);
        assert(operatorCap.check(), message:  "cannot borrow operator capability")

        // create SharedCapabilities resource with withdraw and operator capabilities, and create distribution with it
        issuerRef.createDist(
            sharedCap: <- PDS.createSharedCapabilities(withdrawCap: withdrawCap, operatorCap: operatorCap),
            title: title,
            metadata: metadata,
        )
    }
}
