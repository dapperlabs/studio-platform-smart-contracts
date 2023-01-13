package test

import (
	emulator "github.com/onflow/flow-emulator"
	"testing"

	"github.com/stretchr/testify/assert"
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
