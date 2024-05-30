package test

import (
	"testing"

	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	distributionTitle = "TestDistribution"
)

// Create all required resources for different accounts
func TestMintExampleNFTs(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	// Set exampleNFT (issuer) account and pds account
	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	pdsAccountKey, pdsSigner := accountKeys.NewWithSigner()

	nftAddress, metadataAddress, exampleNFTAddress, _, _ := deployPDSContracts(t, b, exampleNFTAccountKey, exampleNFTSigner, exampleNFTAccountKey, pdsAccountKey,
		pdsSigner)

	t.Run("Should be able to mint a token", func(t *testing.T) {
		// Mint a single NFT with standard royalty cuts and metadata
		mintExampleNFT(t, b,
			accountKeys,
			nftAddress, metadataAddress, exampleNFTAddress,
			exampleNFTAccountKey,
			exampleNFTSigner)

		// Check that the NFT has been minted
		assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
			exampleNFTAddress,
			1,
		)
	})
}

func TestCreatePackIssuer(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	// Set exampleNFT (issuer) account and pds account
	// exampleNFTAccount deploys both the ExampleNFT contract and the PackNFT contract
	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	pdsAccountKey, pdsSigner := accountKeys.NewWithSigner()

	nftAddress, metadataAddress, exampleNFTAddress, _, pdsAddress := deployPDSContracts(
		t,
		b,
		exampleNFTAccountKey,
		exampleNFTSigner,
		exampleNFTAccountKey,
		pdsAccountKey,
		pdsSigner)

	// Mint a single NFT with standard royalty cuts and metadata
	mintExampleNFT(t, b,
		accountKeys,
		nftAddress, metadataAddress, exampleNFTAddress,
		exampleNFTAccountKey,
		exampleNFTSigner)

	// Check that the NFT has been minted
	assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
		exampleNFTAddress,
		1,
	)

	t.Run("Should be able to create a pack issuer NFT", func(t *testing.T) {

		// Assumes issuer is deployer of exampleNFT
		tx := createTxWithTemplateAndAuthorizer(b,
			templates.GenerateCreatePackIssuerTx(pdsAddress),
			exampleNFTAddress,
		)

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
	})
}

func TestCreateDistribution(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	// Set exampleNFT (issuer) account and pds account
	// exampleNFTAccount deploys both the ExampleNFT contract and the PackNFT contract
	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	pdsAccountKey, pdsSigner := accountKeys.NewWithSigner()

	nftAddress, metadataAddress, exampleNFTAddress, iPackNFTAddress, pdsAddress := deployPDSContracts(
		t,
		b,
		exampleNFTAccountKey,
		exampleNFTSigner,
		exampleNFTAccountKey,
		pdsAccountKey,
		pdsSigner)

	// Mint a single NFT with standard royalty cuts and metadata
	mintExampleNFT(t, b,
		accountKeys,
		nftAddress, metadataAddress, exampleNFTAddress,
		exampleNFTAccountKey,
		exampleNFTSigner)

	// Check that the NFT has been minted
	assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
		exampleNFTAddress,
		1,
	)

	// Assumes issuer is deployer of exampleNFT
	// CreatePackIssuerTx
	tx := createTxWithTemplateAndAuthorizer(b,
		templates.GenerateCreatePackIssuerTx(pdsAddress),
		exampleNFTAddress,
	)

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

	// t.Run("Should be able to link NFT provider capability", func(t *testing.T) {

	// 	// Assumes issuer is deployer of exampleNFT
	// 	script = templates.GenerateLinkExampleNFTProviderCapTx(nftAddress, exampleNFTAddress)
	// 	tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
	// 	// Set argument: NFT provider path
	// 	tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})

	// 	serviceSigner, _ := b.ServiceKey().Signer()

	// 	signAndSubmit(
	// 		t, b, tx,
	// 		[]flow.Address{
	// 			b.ServiceKey().Address,
	// 			exampleNFTAddress,
	// 		},
	// 		[]crypto.Signer{
	// 			serviceSigner,
	// 			exampleNFTSigner,
	// 		},
	// 		false,
	// 	)
	// })

	t.Run("Should be able to set pack issuer capability", func(t *testing.T) {

		// Assumes issuer is deployer of exampleNFT
		tx := createTxWithTemplateAndAuthorizer(b,
			templates.GenerateSetPackIssuerCapTx(pdsAddress),
			pdsAddress,
		)
		// Set argument: issuer address
		tx.AddArgument(cadence.NewAddress(exampleNFTAddress))

		serviceSigner, _ := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{
				b.ServiceKey().Address,
				pdsAddress,
			},
			[]crypto.Signer{
				serviceSigner,
				pdsSigner,
			},
			false,
		)
	})

	t.Run("Should be able to create a distribution", func(t *testing.T) {

		// Assumes issuer is deployer of exampleNFT
		tx := createTxWithTemplateAndAuthorizer(b,
			templates.GenerateCreateDistributionTx(pdsAddress, exampleNFTAddress, iPackNFTAddress, nftAddress),
			exampleNFTAddress,
		)
		// Set argument: issuer address
		// tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})
		tx.AddArgument(cadence.String(distributionTitle))
		metadata := []cadence.KeyValuePair{{Key: cadence.String("TestKey"), Value: cadence.String("TestValue")}}
		tx.AddArgument(cadence.NewDictionary(metadata))

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

		// Check that the distribution has been created and has the expected title
		title := executeScriptAndCheck(t, b,
			templates.GenerateGetDistTitleScript(pdsAddress),
			[][]byte{jsoncdc.MustEncode(cadence.UInt64(1))},
		)
		assert.Equal(t, cadence.String(distributionTitle), title)
	})
}

func TestMintPackNFTs(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	// Set exampleNFT (issuer) account and pds account
	// exampleNFTAccount deploys both the ExampleNFT contract and the PackNFT contract
	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	pdsAccountKey, pdsSigner := accountKeys.NewWithSigner()

	nftAddress, metadataAddress, exampleNFTAddress, iPackNFTAddress, pdsAddress := deployPDSContracts(
		t,
		b,
		exampleNFTAccountKey,
		exampleNFTSigner,
		exampleNFTAccountKey,
		pdsAccountKey,
		pdsSigner)

	// Mint a single NFT with standard royalty cuts and metadata
	mintExampleNFT(t, b,
		accountKeys,
		nftAddress, metadataAddress, exampleNFTAddress,
		exampleNFTAccountKey,
		exampleNFTSigner)

	// Check that the NFT has been minted
	assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
		exampleNFTAddress,
		1,
	)

	// Assumes issuer is deployer of exampleNFT
	// CreatePackIssuerTx
	tx := createTxWithTemplateAndAuthorizer(b,
		templates.GenerateCreatePackIssuerTx(pdsAddress),
		exampleNFTAddress,
	)

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

	// // Assumes issuer is deployer of exampleNFT
	// script = templates.GenerateLinkExampleNFTProviderCapTx(nftAddress, exampleNFTAddress)
	// tx = createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
	// // Set argument: NFT provider path
	// // tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})

	// signAndSubmit(
	// 	t, b, tx,
	// 	[]flow.Address{
	// 		b.ServiceKey().Address,
	// 		exampleNFTAddress,
	// 	},
	// 	[]crypto.Signer{
	// 		serviceSigner,
	// 		exampleNFTSigner,
	// 	},
	// 	false,
	// )

	// Assumes issuer is deployer of exampleNFT
	tx = createTxWithTemplateAndAuthorizer(b,
		templates.GenerateSetPackIssuerCapTx(pdsAddress),
		pdsAddress,
	)
	// Set argument: issuer address
	tx.AddArgument(cadence.NewAddress(exampleNFTAddress))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{
			b.ServiceKey().Address,
			pdsAddress,
		},
		[]crypto.Signer{
			serviceSigner,
			pdsSigner,
		},
		false,
	)

	// Assumes issuer is deployer of exampleNFT
	tx = createTxWithTemplateAndAuthorizer(b,
		templates.GenerateCreateDistributionTx(pdsAddress, exampleNFTAddress, iPackNFTAddress, nftAddress),
		exampleNFTAddress,
	)
	// Set argument: issuer address
	// tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})
	tx.AddArgument(cadence.String(distributionTitle))
	metadata := []cadence.KeyValuePair{{Key: cadence.String("TestKey"), Value: cadence.String("TestValue")}}
	tx.AddArgument(cadence.NewDictionary(metadata))

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

	title := executeScriptAndCheck(t, b,
		templates.GenerateGetDistTitleScript(pdsAddress),
		[][]byte{jsoncdc.MustEncode(cadence.UInt64(1))},
	)
	assert.Equal(t, cadence.String(distributionTitle), title)

	t.Run("Should be able to mint a pack NFT", func(t *testing.T) {
		// Assumes issuer is deployer of exampleNFT
		tx := createTxWithTemplateAndAuthorizer(b,
			templates.GenerateMintPackNFTTx(pdsAddress, exampleNFTAddress, nftAddress),
			pdsAddress,
		)
		// Set argument: issuer address
		tx.AddArgument(cadence.UInt64(1))
		tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.String("4b82931fe40fce9b60b68207171dde5d4f07f070e669baf7e08cee18d27c5ef8")}))
		tx.AddArgument(cadence.NewAddress(exampleNFTAddress))

		serviceSigner, _ := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{
				b.ServiceKey().Address,
				pdsAddress,
			},
			[]crypto.Signer{
				serviceSigner,
				pdsSigner,
			},
			false,
		)

		assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
			exampleNFTAddress,
			1,
		)
	})
}

func TestOpenPackNFT(t *testing.T) {
	b, accountKeys := newTestSetup(t)

	// Set exampleNFT (issuer) account and pds account
	// exampleNFTAccount deploys both the ExampleNFT contract and the PackNFT contract
	exampleNFTAccountKey, exampleNFTSigner := accountKeys.NewWithSigner()
	pdsAccountKey, pdsSigner := accountKeys.NewWithSigner()

	nftAddress, metadataAddress, exampleNFTAddress, iPackNFTAddress, pdsAddress := deployPDSContracts(
		t,
		b,
		exampleNFTAccountKey,
		exampleNFTSigner,
		exampleNFTAccountKey,
		pdsAccountKey,
		pdsSigner)

	// Mint a single NFT with standard royalty cuts and metadata
	mintExampleNFT(t, b,
		accountKeys,
		nftAddress, metadataAddress, exampleNFTAddress,
		exampleNFTAccountKey,
		exampleNFTSigner)

	// Check that the NFT has been minted
	assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
		exampleNFTAddress,
		1,
	)

	// Assumes issuer is deployer of exampleNFT
	// CreatePackIssuerTx
	tx := createTxWithTemplateAndAuthorizer(b,
		templates.GenerateCreatePackIssuerTx(pdsAddress),
		exampleNFTAddress,
	)

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

	// // Assumes issuer is deployer of exampleNFT
	// script = templates.GenerateLinkExampleNFTProviderCapTx(nftAddress, exampleNFTAddress)
	// tx = createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
	// // Set argument: NFT provider path
	// // tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})

	// signAndSubmit(
	// 	t, b, tx,
	// 	[]flow.Address{
	// 		b.ServiceKey().Address,
	// 		exampleNFTAddress,
	// 	},
	// 	[]crypto.Signer{
	// 		serviceSigner,
	// 		exampleNFTSigner,
	// 	},
	// 	false,
	// )

	// Assumes issuer is deployer of exampleNFT
	tx = createTxWithTemplateAndAuthorizer(b,
		templates.GenerateSetPackIssuerCapTx(pdsAddress),
		pdsAddress,
	)
	// Set argument: issuer address
	tx.AddArgument(cadence.NewAddress(exampleNFTAddress))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{
			b.ServiceKey().Address,
			pdsAddress,
		},
		[]crypto.Signer{
			serviceSigner,
			pdsSigner,
		},
		false,
	)

	// Assumes issuer is deployer of exampleNFT
	tx = createTxWithTemplateAndAuthorizer(b,
		templates.GenerateCreateDistributionTx(pdsAddress, exampleNFTAddress, iPackNFTAddress, nftAddress),
		exampleNFTAddress,
	)
	// Set argument: issuer address
	tx.AddArgument(cadence.String(distributionTitle))
	metadata := []cadence.KeyValuePair{{Key: cadence.String("TestKey"), Value: cadence.String("TestValue")}}
	tx.AddArgument(cadence.NewDictionary(metadata))

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

	title := executeScriptAndCheck(t, b,
		templates.GenerateGetDistTitleScript(pdsAddress),
		[][]byte{jsoncdc.MustEncode(cadence.UInt64(1))},
	)
	assert.Equal(t, cadence.String(distributionTitle), title)

	t.Run("Should be able to mint a pack NFT", func(t *testing.T) {
		// Assumes issuer is deployer of exampleNFT
		tx := createTxWithTemplateAndAuthorizer(b,
			templates.GenerateMintPackNFTTx(pdsAddress, exampleNFTAddress, nftAddress),
			pdsAddress,
		)
		// Set argument: issuer address
		tx.AddArgument(cadence.UInt64(1))
		tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.String("4b82931fe40fce9b60b68207171dde5d4f07f070e669baf7e08cee18d27c5ef8")}))
		tx.AddArgument(cadence.NewAddress(exampleNFTAddress))

		serviceSigner, _ := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{
				b.ServiceKey().Address,
				pdsAddress,
			},
			[]crypto.Signer{
				serviceSigner,
				pdsSigner,
			},
			false,
		)

		assertCollectionLength(t, b, nftAddress, exampleNFTAddress, metadataAddress,
			exampleNFTAddress,
			1,
		)
	})
}
