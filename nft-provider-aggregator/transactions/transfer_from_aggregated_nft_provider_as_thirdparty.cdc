import NonFungibleToken from "NonFungibleToken"
import ExampleNFT from "ExampleNFT"
import NFTProviderAggregator from "NFTProviderAggregator"

/// Transaction signed by a manager account to transfer a NFT from aggregated NFT provider
/// to the recipient's collection.
///
/// @param recipient: The recipient address of the NFT to transfer.
/// @param withdrawID: The ID of the NFT to withdraw.
///
transaction(recipient: Address, withdrawID: UInt64) {

    let depositRef: &{NonFungibleToken.CollectionPublic}
    let aggregatedProviderRef: auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}

    prepare(
        signer: auth(BorrowValue, Storage) &Account,
    ) {
        // Get recipient account
        let recipient = getAccount(recipient)

        // Borrow a public reference to the receivers collection
        self.depositRef = recipient
            .capabilities.borrow<&{NonFungibleToken.CollectionPublic}>(ExampleNFT.CollectionPublicPath)!

        // Borrow a reference to the signer's aggregated provider
        // Note: The capability is a claimed capability stored in StoragePath instead of the usual
        // CapabilityPath.
        self.aggregatedProviderRef = signer.storage.load<
        Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>>(
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
