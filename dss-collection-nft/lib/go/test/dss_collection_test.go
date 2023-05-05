package test

import (
	"strconv"
	"testing"

	emulator "github.com/onflow/flow-emulator"
	"github.com/stretchr/testify/assert"
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------
func TestDSSCollectionDeployContracts(t *testing.T) {
	b := newEmulator()
	DSSCollectionDeployContracts(t, b)
}

func TestDSSCollectionSetupAccount(t *testing.T) {
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
		)
	})
}

func testCreateCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
) {
	productName := "NBA Top Shot"
	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars",
		productName,
		map[string]string{"key": "META"},
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, collectionGroupId)
		assert.Equal(t, collectionGroupId, collectionGroup.ID)
		assert.Equal(t, collectionGroupName, collectionGroup.Name)
		assert.Equal(t, productName, collectionGroup.ProductName)
		assert.Equal(t, true, collectionGroup.Active)
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
		)
	})
}

func testCreateTimeBoundCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
) {
	collectionGroupId := createTimeBoundCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars Description",
		"NBA Top Shot",
		2368296360_00000000,
		map[string]string{"key": "META"},
	)
	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, collectionGroupId)
		assert.Equal(t, collectionGroupId, collectionGroup.ID)
		assert.Equal(t, collectionGroupName, collectionGroup.Name)
		assert.Equal(t, true, collectionGroup.Active)
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

	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		"Top Shot All Stars",
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupId,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, collectionGroupId)
		assert.Equal(t, collectionGroupId, collectionGroup.ID)
		assert.Equal(t, false, collectionGroup.Active)
	}
}

func TestCreateSlot(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to create a new slot", func(t *testing.T) {

		testCreateSlot(
			t,
			b,
			contracts,
			false,
			"OR",
		)
	})
}

func testCreateSlot(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	logicalOperator string,
) {
	collectionGroupID := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		"NBA All Stars",
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)
	required := true

	slotID := createSlot(
		t,
		b,
		contracts,
		false,
		collectionGroupID,
		logicalOperator,
		required,
		map[string]string{"key": "META"},
	)

	if !shouldRevert {
		slot := getSlotData(t, b, contracts, slotID)
		assert.Equal(t, slotID, slot.ID)
		assert.Equal(t, logicalOperator, slot.LogicalOperator)
		assert.Equal(t, required, slot.Required)
	}
}

func TestCreateItemInSlot(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to create a new item", func(t *testing.T) {

		testCreateItemInSlot(
			t,
			b,
			contracts,
			false,
			"100",
			10,
			"edition.id",
		)
	})
}

func testCreateItemInSlot(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	itemID string,
	points uint64,
	itemType string,
) {
	collectionGroupID := createCollectionGroup(
		t,
		b,
		contracts,
		shouldRevert,
		"NBA All Stars",
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)

	slotID := createSlot(
		t,
		b,
		contracts,
		shouldRevert,
		collectionGroupID,
		"OR",
		true,
		map[string]string{"key": "META"},
	)

	comparator := "="
	createItemInSlot(
		t,
		b,
		contracts,
		shouldRevert,
		itemID,
		points,
		itemType,
		comparator,
		slotID,
	)

	if !shouldRevert {
		slot := getSlotData(t, b, contracts, slotID)
		assert.Equal(t, slotID, slot.ID)
		assert.Equal(t, 1, len(slot.Items))
		assert.Equal(t, itemID, slot.Items[0].ItemID)
		assert.Equal(t, points, slot.Items[0].Points)
		assert.Equal(t, itemType, slot.Items[0].ItemType)
		assert.Equal(t, comparator, slot.Items[0].Comparator)
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
	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupId,
	)

	userAddress, userSigner := createAccount(t, b)
	setupDSSCollectionAccount(t, b, userAddress, userSigner, contracts)

	nftLevel := 5

	nftID := mintNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
		collectionGroupId,
		userAddress.String(),
		uint8(nftLevel),
	)

	if !shouldRevert {
		nft := getNFTData(t, b, contracts, userAddress.String(), int(nftID))
		assert.Equal(t, uint64(nftID), nft.ID)
		assert.Equal(t, collectionGroupId, nft.CollectionGroupID)
		assert.Equal(t, userAddress.String(), nft.CompletionAddress)
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

func TestOwnershipCheck(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to determine if a account still owns the NFTs that completed a CG", func(t *testing.T) {
		testOwnershipCheck(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testOwnershipCheck(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	collectionGroupName := "Top Shot All Stars"
	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupId,
	)

	userAddress, userSigner := createAccount(t, b)
	setupDSSCollectionAccount(t, b, userAddress, userSigner, contracts)
	setupExampleNFT(t, b, userAddress, userSigner, contracts)

	user2Address, user2Signer := createAccount(t, b)
	setupExampleNFT(t, b, user2Address, user2Signer, contracts)

	nftLevel := 5

	exampleNftID := mintExampleNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
	)

	mintNFTAndRecordCompletedWith(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
		collectionGroupId,
		userAddress.String(),
		uint8(nftLevel),
		[]uint64{exampleNftID},
	)

	if !shouldRevert {
		ownershipCheck := checkCollectionOwnership(
			t,
			b,
			contracts,
			false,
			userAddress.String(),
			collectionGroupId,
		)

		assert.True(t, ownershipCheck)

		// Transfer NFT to 0x0 address before checking ownership
		transferNFT(
			t,
			b,
			contracts,
			false,
			userAddress,
			user2Signer,
			user2Address.String(),
			uint64(exampleNftID),
		)

		// Ownership void after NFT has been moved
		ownershipCheck2 := checkCollectionOwnership(
			t,
			b,
			contracts,
			false,
			userAddress.String(),
			collectionGroupId,
		)

		assert.False(t, ownershipCheck2)
	}
}

func TestCompletionCount(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to return the number of completion NFTs minted for a given collectionId", func(t *testing.T) {
		testCompletionCount(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testCompletionCount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	collectionGroupName := "Top Shot All Stars"
	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupId,
	)

	userAddress, userSigner := createAccount(t, b)
	setupDSSCollectionAccount(t, b, userAddress, userSigner, contracts)

	user2Address, user2Signer := createAccount(t, b)
	setupDSSCollectionAccount(t, b, user2Address, user2Signer, contracts)

	nftLevel := 5

	mintNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
		collectionGroupId,
		userAddress.String(),
		uint8(nftLevel),
	)

	mintNFT(
		t,
		b,
		contracts,
		false,
		user2Address.String(),
		collectionGroupId,
		user2Address.String(),
		uint8(nftLevel),
	)

	if !shouldRevert {
		expectedNFTCount := uint64(2)
		nftCount := getCollectionGroupNFTCount(t, b, contracts, collectionGroupId)

		assert.Equal(t, expectedNFTCount, nftCount)
	}
}

func TestDuplicateRewardGuard(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should not be able to mint a reward token to the same address for a given collectionId", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				assert.Fail(t, "Expected error due to duplicate reward minting")
			}
		}()
		testDuplicateRewardGuard(
			t,
			b,
			contracts,
			true,
		)
	})
}

func testDuplicateRewardGuard(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	collectionGroupName := "Top Shot All Stars"
	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars",
		"NBA Top Shot",
		map[string]string{"key": "META"},
	)

	closeCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupId,
	)

	userAddress, userSigner := createAccount(t, b)
	setupDSSCollectionAccount(t, b, userAddress, userSigner, contracts)

	nftLevel := 5

	mintNFTAndRecordCompletedWith(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
		collectionGroupId,
		userAddress.String(),
		uint8(nftLevel),
		[]uint64{0},
	)

	// Will panic since userAddress already minted above
	mintNFT(
		t,
		b,
		contracts,
		true,
		userAddress.String(),
		collectionGroupId,
		userAddress.String(),
		uint8(nftLevel),
	)

}
