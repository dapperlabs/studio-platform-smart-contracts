package test

import (
	emulator "github.com/onflow/flow-emulator"
	"github.com/stretchr/testify/assert"
	"strconv"
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
			"A.0xf8d6e0586b0a20c7.NFT",
		)
	})
}

func testCreateCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	typeName string,
) {
	createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		typeName,
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
			"A.0xf8d6e0586b0a20c7.NFT",
		)
	})
}

func testCreateTimeBoundCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	typeName string,
) {
	createTimeBoundCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		typeName,
		1673986190,
		2368296360,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, 1)
		assert.Equal(t, uint64(1), collectionGroup.ID)
		assert.Equal(t, collectionGroupName, collectionGroup.Name)
		assert.Equal(t, typeName, collectionGroup.TypeName)
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
		"A.0xf8d6e0586b0a20c7.NFT",
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

//func TestAddNFTToCollectionGroup(t *testing.T) {
//	b := newEmulator()
//	contracts := DSSCollectionDeployContracts(t, b)
//	t.Run("Should be able to add an NFT to an open collection group", func(t *testing.T) {
//		testAddNFTToCollectionGroup(
//			t,
//			b,
//			contracts,
//			false,
//		)
//	})
//}
//
//func testAddNFTToCollectionGroup(
//	t *testing.T,
//	b *emulator.Blockchain,
//	contracts Contracts,
//	shouldRevert bool,
//) {
//
//	createCollectionGroup(
//		t,
//		b,
//		contracts,
//		false,
//		"TopDunkers",
//		"A.0xf8d6e0586b0a20c7.NFT",
//	)
//
//	addNFTToCollectionGroup(
//		t,
//		b,
//		contracts,
//		false,
//		1,
//		100,
//	)
//
//	if !shouldRevert {
//		collectionGroup := getCollectionGroupData(t, b, contracts, 1)
//		assert.Equal(t, uint64(1), collectionGroup.ID)
//		assert.Equal(t, true, collectionGroup.Open)
//		//assert.Equal(t, true, collectionGroup.NFTIDInCollectionGroup[100])
//		//assert.Equal(t, 1, len(collectionGroup.NFTIDInCollectionGroup))
//	}
//}

func TestMintNFT(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to mint an nft", func(t *testing.T) {
		testMintNFT(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testMintNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	collectionGroupName := "Top Shot All Stars"
	createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"A.0xf8d6e0586b0a20c7.NFT",
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

	nftLevel := 5

	mintNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
		1,
		userAddress.String(),
		uint64(nftLevel),
	)

	if !shouldRevert {
		nftID := 1
		nft := getNFTData(t, b, contracts, userAddress.String(), nftID)
		assert.Equal(t, uint64(nftID), nft.ID)
		assert.Equal(t, uint64(nftID), nft.CollectionGroupID)
		assert.Equal(t, userAddress.String(), nft.CompletedBy)
		assert.NotNil(t, nft.CompletionDate)

		displayView := getDSSCollectionNFTDisplayMetadataView(
			t,
			b,
			contracts,
			userAddress,
			uint64(nftID),
		)
		assert.Contains(t, displayView.Name, collectionGroupName)
		assert.Contains(t, displayView.Name, strconv.Itoa(nftLevel))
		assert.Contains(t, displayView.Description, userAddress.String())
		assert.NotNil(t, displayView.ImageURL)
	}
}
