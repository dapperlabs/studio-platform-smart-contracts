import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import MetadataViews from "MetadataViews"

transaction (nftWithdrawCapPath: StoragePath) {

    prepare(signer: auth(Storage, Capabilities) &Account) {
        let collectionData = ExampleNFT.resolveContractView(resourceType: nil, viewType: Type<MetadataViews.NFTCollectionData>()) as! MetadataViews.NFTCollectionData?
            ?? panic("ViewResolver does not resolve NFTCollectionData view")

        // Return early if the account already has a collection
        if signer.storage.borrow<&ExampleNFT.Collection>(from: collectionData.storagePath) != nil {
            return
        }

        // Create and save collection to account storage
        signer.storage.save(
            <- ExampleNFT.createEmptyCollection(nftType: Type<@ExampleNFT.NFT>()),
            to: collectionData.storagePath,
        )

        // Create and save authorized collection capability to account storage
        let withdrawCap = signer.capabilities.storage.issue<auth(NonFungibleToken.Withdraw) &ExampleNFT.Collection>(collectionData.storagePath)
        signer.capabilities.storage.getController(byCapabilityID: withdrawCap.id)!.setTag("PDSwithdrawCap")
        signer.storage.save(
            withdrawCap,
            to: nftWithdrawCapPath
        )

        // create a public capability for the collection
        signer.capabilities.unpublish(collectionData.publicPath)
        let collectionCap = signer.capabilities.storage.issue<&ExampleNFT.Collection>(collectionData.storagePath)
        signer.capabilities.publish(collectionCap, at: collectionData.publicPath)
    }
}
