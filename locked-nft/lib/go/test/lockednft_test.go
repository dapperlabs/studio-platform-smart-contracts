package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------
func TestDSSCollectionDeployContracts(t *testing.T) {
	b := newEmulator()
	LockedNFTDeployContracts(t, b)
}

func TestDSSCollectionSetupAccount(t *testing.T) {
	b := newEmulator()
	contracts := LockedNFTDeployContracts(t, b)
	userAddress, userSigner := createAccount(t, b)
	setupLockedNFTAccount(t, b, userAddress, userSigner, contracts)

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
