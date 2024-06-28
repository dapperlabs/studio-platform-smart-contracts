import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import ExampleNFT from "../contracts/ExampleNFT.cdc"
import NFTProviderAggregator from "../contracts/NFTProviderAggregator.cdc"

/// Transaction signed by a manager account to transfer a NFT from aggregated NFT provider
/// to the recipient's collection.
///
/// @param recipient: The recipient address of the NFT to transfer.
/// @param withdrawID: The ID of the NFT to withdraw.
///
transaction(recipient: Address, withdrawID: UInt64) {

    let depositRef: &AnyResource{NonFungibleToken.CollectionPublic}
    let aggregatedProviderRef: &AnyResource{NonFungibleToken.Provider}

    prepare(
        signer: AuthAccount,
    ) {
        // Get recipient account
        let recipient = getAccount(recipient)

        // Borrow a public reference to the receivers collection
        self.depositRef = recipient
            .getCapability(ExampleNFT.CollectionPublicPath)!
            .borrow<&AnyResource{NonFungibleToken.CollectionPublic}>()!

        // Borrow a reference to the signer's aggregated provider
        // Note: The capability is a claimed capability stored in StoragePath instead of the usual
        // CapabilityPath.
        self.aggregatedProviderRef = signer.load<
        Capability<&AnyResource{NonFungibleToken.Provider}>>(
            from: NFTProviderAggregator.AggregatorStoragePath)!.borrow()
            ?? panic("Could not get capability and borrow reference")        
    }

    execute {
        // Withdraw the NFT from the aggregated provider
        let nft <- self.aggregatedProviderRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        self.depositRef.deposit(token: <-nft)
    }
}
