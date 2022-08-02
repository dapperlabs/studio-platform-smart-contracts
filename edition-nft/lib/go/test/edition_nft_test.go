package test

import (
	"fmt"
	"testing"

	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/stretchr/testify/assert"
)

//------------------------------------------------------------
// Setup
//------------------------------------------------------------
func TestAllDaySeasonalDeployContracts(t *testing.T) {
	b := newEmulator()
	AllDaySeasonalDeployContracts(t, b)
}

func TestAllDaySeasonalSetupAccount(t *testing.T) {
	b := newEmulator()
	contracts := AllDaySeasonalDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupAllDaySeasonal(t, b, userAddress, userSigner, contracts)

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

//------------------------------------------------------------
// Edition
//------------------------------------------------------------
func TestEdition(t *testing.T) {
	b := newEmulator()
	contracts := AllDaySeasonalDeployContracts(t, b)
	createTestSeasonalEditions(t, b, contracts)
}

func createTestSeasonalEditions(t *testing.T, b *emulator.Blockchain, contracts Contracts) {
	t.Run("Should be able to create a new edition", func(t *testing.T) {
		testCreateSeasonalEdition(
			t,
			b,
			contracts,
			map[string]string{"test play metadata a": "TEST PLAY METADATA A"},
			1,
			false,
		)
	})

	t.Run("Should be able to create a new edition with an incrementing ID from the first", func(t *testing.T) {
		testCreateSeasonalEdition(
			t,
			b,
			contracts,
			map[string]string{"test play metadata a": "TEST PLAY METADATA A"},
			2,
			false,
		)
	})

	t.Run("Should be able to close a edition", func(t *testing.T) {
		testCloseSeasonalEdition(
			t,
			b,
			contracts,
			2,
			false,
		)
	})
}

func testCreateSeasonalEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	metadata map[string]string,
	shouldBeID uint64,
	shouldRevert bool,
) {
	createEdition(
		t,
		b,
		contracts,
		metadata,
		false,
	)

	if !shouldRevert {
		series := getEditionData(t, b, contracts, shouldBeID)
		assert.Equal(t, shouldBeID, series.ID)
		assert.Equal(t, true, series.Active)
	}
}

func testCloseSeasonalEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	editionID uint64,
	shouldRevert bool,
) {
	wasActive := getEditionData(t, b, contracts, editionID).Active
	closeEdition(
		t,
		b,
		contracts,
		editionID,
		shouldRevert,
	)

	edition := getEditionData(t, b, contracts, editionID)
	assert.Equal(t, editionID, edition.ID)
	if !shouldRevert {
		assert.Equal(t, false, edition.Active)
	} else {
		assert.Equal(t, wasActive, edition.Active)
	}
}

// ------------------------------------------------------------
// MomentNFTs
// ------------------------------------------------------------
func TestSeasonalNFTs(t *testing.T) {
	b := newEmulator()
	contracts := AllDaySeasonalDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupAllDaySeasonal(t, b, userAddress, userSigner, contracts)

	createTestSeasonalEditions(t, b, contracts)

	t.Run("Should be able to mint a new NFT from an edition that has a maxMintSize", func(t *testing.T) {
		testMintSeasonalNFT(
			t,
			b,
			contracts,
			uint64(1),
			userAddress,
			uint64(1),
			false,
		)
	})

	t.Run("Should be able to mint a second new MomentNFT from an edition that has a maxmintSize", func(t *testing.T) {
		testMintSeasonalNFT(
			t,
			b,
			contracts,
			uint64(1),
			userAddress,
			uint64(2),
			false,
		)
	})

	closeEdition(
		t,
		b,
		contracts,
		uint64(1),
		false,
	)

	t.Run("Should not be able to mint an edition that is already closed", func(t *testing.T) {
		testMintSeasonalNFT(
			t,
			b,
			contracts,
			uint64(1),
			userAddress,
			uint64(3),
			true,
		)
	})
}

func testMintSeasonalNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	editionID uint64,
	userAddress flow.Address,
	shouldBeID uint64,
	shouldRevert bool,
) {
	// Make sure the total supply of NFTs is tracked correctly
	previousSupply := getEditionNFTSupply(t, b, contracts)
	fmt.Printf("sss %d \n", previousSupply)

	mintSeasonalNFT(
		t,
		b,
		contracts,
		userAddress,
		editionID,
		shouldRevert,
	)

	newSupply := getEditionNFTSupply(t, b, contracts)
	if !shouldRevert {
		assert.Equal(t, previousSupply+uint64(1), newSupply)

		nftProperties := getEditionNFTProperties(
			t,
			b,
			contracts,
			userAddress,
			shouldBeID,
		)
		assert.Equal(t, shouldBeID, nftProperties.ID)
		assert.Equal(t, editionID, nftProperties.EditionID)
	} else {
		assert.Equal(t, previousSupply, newSupply)
	}
}
