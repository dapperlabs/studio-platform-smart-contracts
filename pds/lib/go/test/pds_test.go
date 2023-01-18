package test

import (
	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPDSDeployment(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	// Create new keys for the NFT contract account
	// and deploy all the NFT contracts
	exampleNFTAccountKey, _ := accountKeys.NewWithSigner()
	nftAddress, _, exampleNFTAddress, _, _ := deployPDSContracts(t, b, exampleNFTAccountKey, exampleNFTAccountKey)

	t.Run("Should have properly initialized fields after deployment", func(t *testing.T) {

		script := templates.GenerateGetTotalSupplyScript(nftAddress, exampleNFTAddress)
		supply := executeScriptAndCheck(t, b, script, nil)
		assert.Equal(t, cadence.NewUInt64(0), supply)

		//assertCollectionLength(t, b, nftAddress, exampleNFTAddress,
		//	exampleNFTAddress,
		//	0,
		//)
	})
}
