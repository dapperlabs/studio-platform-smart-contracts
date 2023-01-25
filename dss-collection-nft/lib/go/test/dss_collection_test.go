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
	collectionGroupId := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		collectionGroupName,
		"All Stars",
		typeName,
	)

	if !shouldRevert {
		collectionGroup := getCollectionGroupData(t, b, contracts, collectionGroupId)
		assert.Equal(t, collectionGroupId, collectionGroup.ID)
		assert.Equal(t, collectionGroupName, collectionGroup.Name)
		assert.Equal(t, true, collectionGroup.Open)
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
		"A.0xf8d6e0586b0a20c7.NFT",
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
		assert.Equal(t, false, collectionGroup.Open)
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
			"A.0xf8d6e0586b0a20c7.NFT",
		)
	})
}

func testCreateSlot(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	logicalOperator string,
	typeName string,
) {
	collectionGroupID := createCollectionGroup(
		t,
		b,
		contracts,
		false,
		"NBA All Stars",
		"All Stars",
		typeName,
	)

	slotID := createSlot(
		t,
		b,
		contracts,
		false,
		collectionGroupID,
		logicalOperator,
		typeName,
	)

	if !shouldRevert {
		slot := getSlotData(t, b, contracts, slotID)
		assert.Equal(t, slotID, slot.ID)
		assert.Equal(t, logicalOperator, slot.LogicalOperator)
		assert.Equal(t, typeName, slot.TypeName)
	}
}

func TestCreateItem(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to create a new item", func(t *testing.T) {

		testCreateItem(
			t,
			b,
			contracts,
			false,
			100,
			10,
			"edition.id",
		)
	})
}

func testCreateItem(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	itemID uint64,
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
		"A.0xf8d6e0586b0a20c7.NFT",
	)

	_ = createSlot(
		t,
		b,
		contracts,
		shouldRevert,
		collectionGroupID,
		"OR",
		"A.0xf8d6e0586b0a20c7.NFT",
	)

	id := createItem(
		t,
		b,
		contracts,
		shouldRevert,
		itemID,
		points,
		itemType,
	)

	if !shouldRevert {
		item := getItemData(t, b, contracts, id)
		assert.Equal(t, id, item.ID)
		assert.Equal(t, itemID, item.ItemID)
		assert.Equal(t, points, item.Points)
		assert.Equal(t, itemType, item.ItemType)
	}
}

func TestAddItemToSlot(t *testing.T) {
	b := newEmulator()
	contracts := DSSCollectionDeployContracts(t, b)
	t.Run("Should be able to add an item to a slot", func(t *testing.T) {
		testAddItemToSlot(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testAddItemToSlot(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {

	collectionGroupID := createCollectionGroup(
		t,
		b,
		contracts,
		shouldRevert,
		"NBA All Stars",
		"All Stars",
		"A.0xf8d6e0586b0a20c7.NFT",
	)

	slotID := createSlot(
		t,
		b,
		contracts,
		shouldRevert,
		collectionGroupID,
		"OR",
		"A.0xf8d6e0586b0a20c7.NFT",
	)

	itemID := 100
	points := 10
	itemType := "edition.id"
	id := createItem(
		t,
		b,
		contracts,
		shouldRevert,
		uint64(itemID),
		uint64(points),
		itemType,
	)

	addItemToSlot(
		t,
		b,
		contracts,
		false,
		slotID,
		id,
	)

	if !shouldRevert {
		slot := getSlotData(t, b, contracts, slotID)
		assert.Equal(t, slotID, slot.ID)
		assert.Equal(t, 1, len(slot.Items))
		assert.Equal(t, uint64(itemID), slot.Items[0].ItemID)
		assert.Equal(t, uint64(points), slot.Items[0].Points)
		assert.Equal(t, itemType, slot.Items[0].ItemType)
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
		"A.0xf8d6e0586b0a20c7.NFT",
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
