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
	NFTLockerDeployContracts(t, b)
}

func TestDSSCollectionSetupAccount(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupNFTLockerAccount(t, b, userAddress, userSigner, contracts)

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

func TestLockNFT(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should be able to mint and lock an nft", func(t *testing.T) {
		testLockNFT(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testLockNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	var duration uint64 = 10
	userAddress, userSigner := createAccount(t, b)
	setupNFTLockerAccount(t, b, userAddress, userSigner, contracts)
	setupExampleNFT(t, b, userAddress, userSigner, contracts)

	exampleNftID := mintExampleNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
	)

	lockedAt, lockedUntil := lockNFT(
		t,
		b,
		contracts,
		false,
		userAddress,
		userSigner,
		exampleNftID,
		duration,
	)

	assert.Equal(t, lockedAt+duration, lockedUntil)
}

func TestUnlockNFT(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should be able to mint and lock an nft", func(t *testing.T) {
		testUnlockNFT(
			t,
			b,
			contracts,
		)
	})
}

func testUnlockNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
) {
	var duration uint64 = 200
	userAddress, userSigner := createAccount(t, b)
	setupNFTLockerAccount(t, b, userAddress, userSigner, contracts)
	setupExampleNFT(t, b, userAddress, userSigner, contracts)

	exampleNftID := mintExampleNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
	)

	lockedAt, lockedUntil := lockNFT(
		t,
		b,
		contracts,
		false,
		userAddress,
		userSigner,
		exampleNftID,
		duration,
	)
	assert.Equal(t, lockedAt+duration, lockedUntil)

	unlockNFT(
		t,
		b,
		contracts,
		true,
		userAddress,
		userSigner,
		exampleNftID,
	)

	//assert.Equal(t, lockedAt+duration, lockedUntil)
}
