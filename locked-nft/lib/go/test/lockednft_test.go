package test

import (
	emulator "github.com/onflow/flow-emulator"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestReLockNFT(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should fail to relock a locked nft", func(t *testing.T) {
		testReLockNFT(
			t,
			b,
			contracts,
		)
	})
}

func testReLockNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
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

	// should revert lock of already locked nft
	lockedAt, lockedUntil = lockNFT(
		t,
		b,
		contracts,
		true,
		userAddress,
		userSigner,
		exampleNftID,
		duration,
	)
}

func TestUnlockNFTPanic(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should be able to mint and throw panic unlock an nft", func(t *testing.T) {
		testUnlockNFTPanic(
			t,
			b,
			contracts,
		)
	})
}

func testUnlockNFTPanic(
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

	assert.Equal(t, lockedAt+duration, lockedUntil)
}

func TestUnlockNFT(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should be able to mint, lock, and unlock an nft", func(t *testing.T) {
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
	var duration uint64 = 0
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

	// fast-forward block time past unlock duration
	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
		b.CommitBlock()
	}

	unlockNFT(
		t,
		b,
		contracts,
		false,
		userAddress,
		userSigner,
		exampleNftID,
	)
}

func TestExtendLock(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should be able to extend the lock of an NFT", func(t *testing.T) {
		testExtendLock(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testExtendLock(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	var duration uint64 = 10
	var extendedDuration uint64 = 1000
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

	lockedDataPre := getLockedTokenData(
		t,
		b,
		contracts,
		exampleNftID,
	)

	extendLock(
		t,
		b,
		contracts,
		false,
		userAddress,
		userSigner,
		exampleNftID,
		extendedDuration,
	)

	lockedDataPost := getLockedTokenData(
		t,
		b,
		contracts,
		exampleNftID,
	)

	assert.Equal(t, lockedAt+duration+extendedDuration, lockedDataPost.LockedUntil)
	assert.Less(t, lockedDataPre.LockedUntil, lockedDataPost.LockedUntil)
}

func TestExtendLockFail(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should fail to extend the lock of an NFT that has not been locked", func(t *testing.T) {
		testExtendLockFail(
			t,
			b,
			contracts,
			true,
		)
	})
}

func testExtendLockFail(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	var extendedDuration uint64 = 1000
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

	extendLock(
		t,
		b,
		contracts,
		shouldRevert,
		userAddress,
		userSigner,
		exampleNftID,
		extendedDuration,
	)
}

func TestSwapLock(t *testing.T) {
	b := newEmulator()
	contracts := NFTLockerDeployContracts(t, b)
	t.Run("Should be able to swap the lock of an unlockable NFT", func(t *testing.T) {
		testSwapLock(
			t,
			b,
			contracts,
			false,
		)
	})
}

func testSwapLock(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	var duration uint64 = 0
	userAddress, userSigner := createAccount(t, b)
	setupNFTLockerAccount(t, b, userAddress, userSigner, contracts)
	setupExampleNFT(t, b, userAddress, userSigner, contracts)

	exampleNft1ID := mintExampleNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
	)

	exampleNft2ID := mintExampleNFT(
		t,
		b,
		contracts,
		false,
		userAddress.String(),
	)

	nftInventoryPre := getNFTInventory(
		t,
		b,
		contracts,
		userAddress,
	)

	assert.Equal(t, 2, len(nftInventoryPre))

	lockedAt, lockedUntil := lockNFT(
		t,
		b,
		contracts,
		false,
		userAddress,
		userSigner,
		exampleNft1ID,
		duration,
	)

	assert.Equal(t, lockedAt+duration, lockedUntil)

	// fast-forward block time past unlock duration
	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
		b.CommitBlock()
	}

	swapLock(
		t,
		b,
		contracts,
		false,
		userAddress,
		userSigner,
		exampleNft1ID,
		exampleNft2ID,
		duration,
	)

	lockedDataNFT2 := getLockedTokenData(
		t,
		b,
		contracts,
		exampleNft2ID,
	)

	assert.Equal(t, lockedDataNFT2.LockedAt+duration, lockedDataNFT2.LockedUntil)

	nftInventoryPost := getNFTInventory(
		t,
		b,
		contracts,
		userAddress,
	)

	assert.Equal(t, true, arrayContains(nftInventoryPost, exampleNft1ID))
	assert.Equal(t, false, arrayContains(nftInventoryPost, exampleNft2ID))
}

func arrayContains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
