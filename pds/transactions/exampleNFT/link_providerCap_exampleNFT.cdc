import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

transaction (nftWithdrawCapPath: StoragePath) {

    prepare(signer: auth(Storage, Capabilities) &Account) {
        // Return early if the account already has an authorized collection capability
        if signer.storage.borrow<auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection>(from: nftWithdrawCapPath) != nil {
            return
        }

        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        // Create and save authorized collection capability to account storage
        let withdrawCap = signer.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection>(collectionData.storagePath)
        signer.capabilities.storage.getController(byCapabilityID: withdrawCap.id)!.setTag("PDSwithdrawCap")
        signer.storage.save(
            withdrawCap,
            to: nftWithdrawCapPath
        )
    }
}
