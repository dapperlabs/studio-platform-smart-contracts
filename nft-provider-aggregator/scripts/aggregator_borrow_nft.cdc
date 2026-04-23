import NFTProviderAggregator from "NFTProviderAggregator"
import NonFungibleToken from "NonFungibleToken"

// Tries to borrow an NFT by ID from the Aggregator resource stored in `address`.
// Returns the NFT's Type identifier if found, nil otherwise.

access(all) fun main(address: Address, nftID: UInt64): &{NonFungibleToken.NFT}? {
    let account = getAuthAccount<auth(BorrowValue) &Account>(address)

    let aggregator = account.storage.borrow<&NFTProviderAggregator.Aggregator>(
        from: NFTProviderAggregator.AggregatorStoragePath
    ) ?? panic("Aggregator resource not found at ".concat(NFTProviderAggregator.AggregatorStoragePath.toString()))

    let nft = aggregator.borrowNFT(id: nftID)
    return nft
}
