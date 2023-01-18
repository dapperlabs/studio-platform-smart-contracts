package test

import (
	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
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

		assertCollectionLength(t, b, nftAddress, exampleNFTAddress,
			exampleNFTAddress,
			0,
		)
	})
}

// Create all required resources for different accounts
func TestMintExampleNFTs(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	nftAddress, metadataAddress, exampleNFTAddress, _, _ := deployPDSContracts(t, b, exampleNFTAccountKey, exampleNFTAccountKey)

	t.Run("Should be able to mint a token", func(t *testing.T) {
		// Mint a single NFT with standard royalty cuts and metadata
		mintExampleNFT(t, b,
			accountKeys,
			nftAddress, metadataAddress, exampleNFTAddress,
			exampleNFTAccountKey,
			exampleNFTSigner)

		script := templates.GenerateBorrowNFTScript(nftAddress, exampleNFTAddress)
		executeScriptAndCheck(
			t, b,
			script,
			[][]byte{
				jsoncdc.MustEncode(cadence.NewAddress(exampleNFTAddress)),
				jsoncdc.MustEncode(cadence.NewUInt64(0)),
			},
		)

		script = templates.GenerateGetTotalSupplyScript(nftAddress, exampleNFTAddress)
		supply := executeScriptAndCheck(t, b, script, nil)
		assert.Equal(t, cadence.NewUInt64(1), supply)

		assertCollectionLength(t, b, nftAddress, exampleNFTAddress,
			exampleNFTAddress,
			1,
		)
	})
}

// Create all required resources for different accounts
func TestMintPackNFTs(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	nftAddress, metadataAddress, exampleNFTAddress, _, _ := deployPDSContracts(t, b, exampleNFTAccountKey, exampleNFTAccountKey)

	t.Run("Should be able to mint a token", func(t *testing.T) {
		// Mint a single NFT with standard royalty cuts and metadata
		mintExampleNFT(t, b,
			accountKeys,
			nftAddress, metadataAddress, exampleNFTAddress,
			exampleNFTAccountKey,
			exampleNFTSigner)

	})
}
