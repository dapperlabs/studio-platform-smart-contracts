package test

import (
	"fmt"
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------
func TestNftProviderAggregatorDeployContracts(t *testing.T) {
	b := newEmulator()
	NFTProviderAggregatorDeployContracts(t, b)
}

const (
	withdrawCapTag = "nft-provider-aggregator"
)

func TestExampleNFTSetupAndMinting(t *testing.T) {
	b := newEmulator()
	contracts := NFTProviderAggregatorDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupExampleNft(t, b, userAddress, userSigner, contracts)

	var nftId uint64

	t.Run("Account should have empty collection", func(t *testing.T) {
		result := getExampleNFTCollectionIds(
			t,
			b,
			contracts,
			userAddress,
		)
		assert.Nil(t, result)
	})

	t.Run("Mint Example NFT", func(t *testing.T) {
		nftId = mintExampleNFT(t, b, contracts, userAddress, false)
	})

	t.Run("Account should have empty collection", func(t *testing.T) {
		result := getExampleNFTCollectionIds(
			t,
			b,
			contracts,
			userAddress,
		)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, nftId, result[0].(uint64))
	})
}

func TestNFTProviderAggregator(t *testing.T) {
	b := newEmulator()
	contracts := NFTProviderAggregatorDeployContracts(t, b)

	// Define constants
	const (
		nftWithdrawCapStoragePathID = "NFT_PROVIDER_CAP"
		nftCollectionStoragePathID  = "exampleNFTCollection"
		capabilityPublicationID1    = "1"
		capabilityPublicationID2    = "2"
		capabilityPublicationID3    = "3"
	)

	var (
		supplier1NftId, supplier2NftId                 uint64
		managerCollectionUUID, supplier2CollectionUUID uint64
	)

	// Setup accounts and mint NFTs
	// setupExampleNft(t, b, contracts.NFTProviderAggregatorAddress, contracts.NFTProviderAggregatorSigner, contracts)
	supplier1Address, supplier1Signer := createAccount(t, b)
	setupExampleNft(t, b, supplier1Address, supplier1Signer, contracts)
	supplier1NftId = mintExampleNFT(t, b, contracts, supplier1Address, false)
	t.Log("Supplier 1 NFT ID:", supplier1NftId)
	supplier2Address, supplier2Signer := createAccount(t, b)
	setupExampleNft(t, b, supplier2Address, supplier2Signer, contracts)
	supplier2NftId = mintExampleNFT(t, b, contracts, supplier2Address, false)
	t.Log("Supplier 2 NFT ID:", supplier2NftId)
	thirdPartyAddress, thirdPartySigner := createAccount(t, b)
	setupExampleNft(t, b, thirdPartyAddress, thirdPartySigner, contracts)

	t.Run("should be able to bootstrap an Aggregator resource", func(t *testing.T) {
		bootstrapAggregatorResource(t, b, contracts,
			fmt.Sprintf("A.%s.ExampleNFT.Collection", contracts.NFTProviderAggregatorAddress.String()),
			[]flow.Address{supplier1Address, supplier2Address},
			[]string{capabilityPublicationID1, capabilityPublicationID2},
			false,
		)
	})

	t.Run("should be able to add a NFT provider capability as manager", func(t *testing.T) {
		collectionUUIDs := addNftWithdrawCapAsManager(t, b, contracts,
			nftWithdrawCapStoragePathID,
			nftCollectionStoragePathID,
			withdrawCapTag,
			false,
		)
		assert.Len(t, collectionUUIDs, 1)
		managerCollectionUUID = collectionUUIDs[0]
	})

	t.Run("should NOT be able to add a NFT provider capability that is already existing as manager", func(t *testing.T) {
		addNftWithdrawCapAsManager(t, b, contracts,
			nftWithdrawCapStoragePathID,
			nftCollectionStoragePathID,
			withdrawCapTag,
			true,
		)
	})

	t.Run("should be able to bootstrap a Supplier resources", func(t *testing.T) {
		bootstrapSupplierResource(t, b, contracts,
			capabilityPublicationID1,
			supplier1Address,
			supplier1Signer,
			false,
		)
	})

	t.Run("should NOT be able to add a NFT provider capability that targets a collection with invalid NFT type", func(t *testing.T) {
		addNftWithdrawCapAsSupplier(t, b, contracts,
			nftWithdrawCapStoragePathID,
			"invalidStoragePath",
			withdrawCapTag,
			supplier1Address,
			supplier1Signer,
			true,
		)
	})

	t.Run("should be able to add a NFT provider capability as supplier", func(t *testing.T) {
		addNftWithdrawCapAsSupplier(t, b, contracts,
			nftWithdrawCapStoragePathID,
			nftCollectionStoragePathID,
			withdrawCapTag,
			supplier1Address,
			supplier1Signer,
			false,
		)
	})

	t.Run("should NOT be able to add a NFT provider capability that is already existing as supplier", func(t *testing.T) {
		addNftWithdrawCapAsSupplier(t, b, contracts,
			nftWithdrawCapStoragePathID,
			nftCollectionStoragePathID,
			withdrawCapTag,
			supplier1Address,
			supplier1Signer,
			true,
		)
	})

	t.Run("should be able to withdraw NFTs from Aggregator's aggregated provider held in both the supplier and the manager's own collections", func(t *testing.T) {
		bootstrapSupplierResource(t, b, contracts,
			capabilityPublicationID2,
			supplier2Address,
			supplier2Signer,
			false,
		)
		uuids := addNftWithdrawCapAsSupplier(t, b, contracts,
			nftWithdrawCapStoragePathID,
			nftCollectionStoragePathID,
			withdrawCapTag,
			supplier2Address,
			supplier2Signer,
			false,
		)
		assert.Len(t, uuids, 1)
		supplier2CollectionUUID = uuids[0]
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			supplier1NftId,
			false,
		)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			supplier2NftId,
			false,
		)
	})

	t.Run("should be able to withdraw NFTs from Aggregator's aggregated provider even if supplier's own collection is empty but manager's is not", func(t *testing.T) {
		nftIDs := getExampleNFTCollectionIds(t, b, contracts, supplier1Address)
		assert.Nil(t, nftIDs)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			supplier1NftId,
			false,
		)
	})

	t.Run("should be able to publish the manager's aggregated NFT provider capability", func(t *testing.T) {
		publishAggregatedNftWithdrawCap(t, b, contracts,
			thirdPartyAddress,
			capabilityPublicationID3,
			false,
		)
	})

	t.Run("should be able to claim the manager's aggregated NFT provider capability and withdraw from it", func(t *testing.T) {
		nftId := mintExampleNFT(t, b, contracts, contracts.NFTProviderAggregatorAddress, false)
		claimAggregatedNftWithdrawCap(t, b, contracts,
			capabilityPublicationID3,
			thirdPartyAddress,
			thirdPartySigner,
			false,
		)
		transferFromAggregatedNftProviderAsThirdParty(t, b, contracts,
			thirdPartyAddress,
			nftId,
			thirdPartyAddress,
			thirdPartySigner,
			false,
		)
	})

	t.Run("should NOT be able to remove a NFT provider capability added by a separate supplier as supplier", func(t *testing.T) {
		removeNftWithdrawCapAsSupplier(t, b, contracts,
			supplier2CollectionUUID,
			supplier1Address,
			supplier1Signer,
			true,
		)
	})

	t.Run("should be able to remove a NFT provider capability added by supplier as manager", func(t *testing.T) {
		removeNftWithdrawCapAsManager(t, b, contracts,
			managerCollectionUUID,
			false,
		)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			supplier1NftId,
			true,
		)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			supplier2NftId,
			true,
		)
	})

	t.Run("should be able to remove supplied NFT provider capabilities when a Supplier resource is destroyed", func(t *testing.T) {
		destroySupplier(t, b, contracts,
			supplier1Address,
			supplier1Signer,
			false,
		)
		nftId := mintExampleNFT(t, b, contracts, supplier1Address, false)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			nftId,
			true,
		)
	})

	t.Run("should be able to withdraw NFTs from Aggregator's aggregated provider even if manager's capability gets revoked but supplier's does not", func(t *testing.T) {
		nftId := mintExampleNFT(t, b, contracts, supplier2Address, false)

		revokeWithdrawCapability(t, b, contracts,
			withdrawCapTag,
			contracts.NFTProviderAggregatorAddress,
			contracts.NFTProviderAggregatorSigner,
			false,
		)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			nftId,
			false,
		)
	})

	t.Run("should NOT be able to withdraw NFTs in manager's collection using Aggregator's aggregated provider if manager has revoked their NFT provider capability", func(t *testing.T) {
		nftId := mintExampleNFT(t, b, contracts, contracts.NFTProviderAggregatorAddress, false)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			nftId,
			true,
		)
	})

	t.Run("should be able to nullify the aggregated NFT provider and child Supplier resources when the Aggregator resource is destroyed", func(t *testing.T) {
		nftId := mintExampleNFT(t, b, contracts, supplier2Address, false)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			nftId,
			false,
		)
		destroyAggregator(t, b, contracts,
			false,
		)
		transferFromAggregatedNftProviderAsManager(t, b, contracts,
			contracts.NFTProviderAggregatorAddress,
			nftId,
			true,
		)
	})
}
