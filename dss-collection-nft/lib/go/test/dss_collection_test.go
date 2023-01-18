package test

import (
	emulator "github.com/onflow/flow-emulator"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------
func TestDSSCollectionDeployContracts(t *testing.T) {
	b := newEmulator()
	DSSCollectionDeployContracts(t, b)
}

func TestEditionNFTSetupAccount(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupDSSCollectionAccount(t, b, userAddress, userSigner, contracts)

	t.Run("Account should be set up", func(t *testing.T) {
		isAccountSetUp := isAccountSetup(
			t,
			b,
			contracts,
			userAddress,
		)
		assert.Equal(t, true, isAccountSetUp)
	})
}

// ------------------------------------------------------------
// Collection Group
// ------------------------------------------------------------
func TestCreateCollectionGroup(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to create a new collection group", func(t *testing.T) {
		testCreateCollectionGroup(
			t,
			b,
			contracts,
			false,
			"Top Shot All Stars",
			"public",
			"tscollection",
		)
	})
}

func testCreateCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	productPathDomain string,
	productPathIdentifier string,
) {
	createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		productPathDomain,
		productPathIdentifier,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, 1)
		assert.Equal(t, uint64(1), collectionGroup.ID)
		assert.Equal(t, collectionGroupName, collectionGroup.Name)
		assert.Equal(t, true, collectionGroup.Open)
		assert.Equal(t, false, collectionGroup.TimeBound)
	}
}

func TestCreateTimeBoundCollectionGroup(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to create a new time-bound collection group", func(t *testing.T) {
		testCreateTimeBoundCollectionGroup(
			t,
			b,
			contracts,
			false,
			"Top Shot All Stars",
			"public",
			"tscollection",
		)
	})
}

func testCreateTimeBoundCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	productPathDomain string,
	productPathIdentifier string,
) {
	createTimeBoundCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		productPathDomain,
		productPathIdentifier,
		1673986190,
		2368296360,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, 1)
		assert.Equal(t, uint64(1), collectionGroup.ID)
		assert.Equal(t, collectionGroupName, collectionGroup.Name)
		assert.Equal(t, true, collectionGroup.Open)
		assert.Equal(t, true, collectionGroup.TimeBound)
	}
}

func TestCloseCollectionGroup(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to close an open collection group", func(t *testing.T) {
		testCloseCollectionGroup(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testCloseCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {

	createCollectionGroup(
		t,
		b,
		contracts,
		false,
		"Top Shot All Stars",
		"public",
		"tscollection",
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		1,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, 1)
		assert.Equal(t, uint64(1), collectionGroup.ID)
		assert.Equal(t, false, collectionGroup.Open)
	}
}

func TestAddNFTToCollectionGroup(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to add an NFT to an open collection group", func(t *testing.T) {
		testAddNFTToCollectionGroup(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testAddNFTToCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {

	createCollectionGroup(
		t,
		b,
		contracts,
		false,
		"TopDunkers",
		"public",
		"tscollection",
	)

	addNFTToCollectionGroup(
		t,
		b,
		contracts,
		false,
		1,
		100,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, 1)
		assert.Equal(t, uint64(1), collectionGroup.ID)
		assert.Equal(t, true, collectionGroup.Open)
		assert.Equal(t, true, collectionGroup.NFTIDInCollectionGroup[100])
		assert.Equal(t, 1, len(collectionGroup.NFTIDInCollectionGroup))
	}
}

func TestMintNFT(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to mint an nft", func(t *testing.T) {
		testMintNFT(
			t,
			b,
			contracts,
			false,
			"0x179b6b1cb6755e31",
		)
	})
}

func testMintNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	recipientAddress string,
) {
	createCollectionGroup(
		t,
		b,
		contracts,
		false,
		"Top Shot All Stars",
		"public",
		"tscollection",
	)

	addNFTToCollectionGroup(
		t,
		b,
		contracts,
		false,
		1,
		100,
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		1,
	)

	userAddress, userSigner := createAccount(t, b)
	setupDSSCollectionAccount(t, b, userAddress, userSigner, contracts)

	mintNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
		1,
		userAddress.String(),
	)

	if !shouldRevert {
		nft := getNFTData(t, b, contracts, userAddress.String(), 1)
		assert.Equal(t, uint64(1), nft.ID)
		assert.Equal(t, uint64(1), nft.CollectionGroupID)
		assert.Equal(t, userAddress.String(), nft.CompletedBy)
		assert.NotNil(t, nft.CompletionDate)
	}
}
