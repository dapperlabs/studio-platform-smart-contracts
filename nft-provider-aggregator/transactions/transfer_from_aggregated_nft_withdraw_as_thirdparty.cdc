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

    let depositRef: &{NonFungibleToken.Collection}
    let aggregatedProviderRef: auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}

    prepare(
        signer: auth(BorrowValue, CopyValue) &Account,
    ) {
        // Get recipient account
        let recipient = getAccount(recipient)

        // Borrow a public reference to the receivers collection
        self.depositRef = recipient
            .capabilities.borrow<&{NonFungibleToken.Collection}>(ExampleNFT.CollectionPublicPath)!

        // Retrieve the aggregated NFT provider's withdraw capability from storage
        let aggregatedNFTWithdrawCapStoragePath = NFTProviderAggregator.convertPrivateToStoragePath(NFTProviderAggregator.AggregatedProviderPrivatePath)
        let aggregatedWithdrawCap: Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}> = signer.storage.copy<
        Capability<auth(NonFungibleToken.Withdraw) &{NonFungibleToken.Provider}>>(
            from: aggregatedNFTWithdrawCapStoragePath)
            ?? panic("Could not retrieve capability from storage!")

        // Borrow a reference to the nft provider aggregator's withdraw capability
        self.aggregatedProviderRef = aggregatedWithdrawCap.borrow()
            ?? panic("Could not borrow reference from capability!")
    }

    execute {
        // Withdraw the NFT from the aggregated provider
        let nft <- self.aggregatedProviderRef.withdraw(withdrawID: withdrawID)

        // Deposit the NFT in the recipient's collection
        self.depositRef.deposit(token: <-nft)
    }
}
