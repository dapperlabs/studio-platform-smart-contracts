package test

import (
	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"
	"testing"
)

//------------------------------------------------------------
// This tests if PDS can support to mint multiple packs
//------------------------------------------------------------
func TestMultiplePackNFTs(t *testing.T) {
	b := newEmulator()
	contracts := PackNftDeployContracts(t, b)
	contracts = PDSDeployContracts(t, b)

	b, accountKeys := newTestSetup(t)

	t.Run("Should be able to create a new edition", func(t *testing.T) {
		// Mint a Pack NFT with standard royalty cuts and metadata
		mintExampleNFT(t, b,
			accountKeys,
			nftAddress, metadataAddress, exampleNFTAddress,
			exampleNFTAccountKey,
			exampleNFTSigner)
	})
}

// Mints a single NFT from the ExampleNFT contract
// with standard metadata fields and royalty cuts
func mintExampleNFT(
	t *testing.T,
	b *emulator.Blockchain,
	accountKeys *test.AccountKeys,
	nftAddress, metadataAddress, exampleNFTAddress flow.Address,
	exampleNFTAccountKey *flow.AccountKey,
	exampleNFTSigner crypto.Signer,
) {

	// Create two new accounts to act as beneficiaries for royalties
	beneficiaryAddress1, _, beneficiarySigner1 := newAccountWithAddress(b, accountKeys)
	setupRoyaltyReceiver(t, b,
		metadataAddress,
		beneficiaryAddress1,
		beneficiarySigner1,
	)
	beneficiaryAddress2, _, beneficiarySigner2 := newAccountWithAddress(b, accountKeys)
	setupRoyaltyReceiver(t, b,
		metadataAddress,
		beneficiaryAddress2,
		beneficiarySigner2,
	)

	// Generate the script that mints a new NFT and deposits it into the recipient's account
	// whose address is the first argument to the transaction
	script := templates.GenerateMintNFTScript(nftAddress, exampleNFTAddress, metadataAddress, flow.HexToAddress(emulatorFTAddress))

	// Create the transaction object with the generated script and authorizer
	tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)

	// Assemble the cut information for royalties
	cut1 := CadenceUFix64("0.25")
	cut2 := CadenceUFix64("0.40")
	cuts := []cadence.Value{cut1, cut2}

	// Assemble the royalty description and beneficiary addresses to get their receivers
	royaltyDescriptions := []cadence.Value{cadence.String("Minter royalty"), cadence.String("Creator royalty")}
	royaltyBeneficiaries := []cadence.Value{cadence.NewAddress(beneficiaryAddress1), cadence.NewAddress(beneficiaryAddress2)}

	// First argument is the recipient of the newly minted NFT
	tx.AddArgument(cadence.NewAddress(exampleNFTAddress))
	tx.AddArgument(cadence.String("Example NFT 0"))
	tx.AddArgument(cadence.String("This is an example NFT"))
	tx.AddArgument(cadence.String("example.jpeg"))
	tx.AddArgument(cadence.NewArray(cuts))
	tx.AddArgument(cadence.NewArray(royaltyDescriptions))
	tx.AddArgument(cadence.NewArray(royaltyBeneficiaries))

	serviceSigner, _ := b.ServiceKey().Signer()

	signAndSubmit(
		t, b, tx,
		[]flow.Address{
			b.ServiceKey().Address,
			exampleNFTAddress,
		},
		[]crypto.Signer{
			serviceSigner,
			exampleNFTSigner,
		},
		false,
	)
}
