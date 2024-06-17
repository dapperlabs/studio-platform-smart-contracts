import PDS from "PDS"
import PackNFT from "PackNFT"
import ExampleNFT from "ExampleNFT"
import IPackNFT from "IPackNFT"
import NonFungibleToken from "NonFungibleToken"
import MetadataViews from "MetadataViews"

transaction(title: String, metadata: {String: String}) {
    prepare (issuer: auth(BorrowValue, Capabilities) &Account) {
        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        let i = issuer.storage.borrow<auth(PDS.CreateDist) &PDS.PackIssuer>(from: PDS.PackIssuerStoragePath)
            ?? panic ("issuer does not have PackIssuer resource")

        // issuer must have a PackNFT collection
        let withdrawCap = issuer.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>(collectionData.storagePath);
        let operatorCap = issuer.capabilities.storage.issue<auth(IPackNFT.Operate) &{IPackNFT.IOperator}>(PackNFT.OperatorStoragePath);
        assert(withdrawCap.check(), message:  "cannot borrow withdraw capability")
        assert(operatorCap.check(), message:  "cannot borrow operator capability")

        let sc <- PDS.createSharedCapabilities ( withdrawCap: withdrawCap, operatorCap: operatorCap )
        i.createDist(sharedCap: <-sc, title: title, metadata: metadata)
    }
}
