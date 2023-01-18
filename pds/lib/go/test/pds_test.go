package test

import (
	"testing"
)

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

		//script := templates.GenerateBorrowNFTScript(nftAddress, exampleNFTAddress)
		//executeScriptAndCheck(
		//	t, b,
		//	script,
		//	[][]byte{
		//		jsoncdc.MustEncode(cadence.NewAddress(exampleNFTAddress)),
		//		jsoncdc.MustEncode(cadence.NewUInt64(0)),
		//	},
		//)

		//script = templates.GenerateGetTotalSupplyScript(nftAddress, exampleNFTAddress)
		//supply := executeScriptAndCheck(t, b, script, nil)
		//assert.Equal(t, cadence.NewUInt64(1), supply)
		//
		//assertCollectionLength(t, b, nftAddress, exampleNFTAddress,
		//	exampleNFTAddress,
		//	1,
		//)
	})
}
