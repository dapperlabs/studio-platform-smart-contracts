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
	script := templates.GenerateBorrowNFTScript(nftAddress, exampleNFTAddress)
	executeScriptAndCheck(
		t, b,
		script,
		[][]byte{
			jsoncdc.MustEncode(cadence.NewAddress(exampleNFTAddress)),
			jsoncdc.MustEncode(cadence.NewUInt64(0)),
		},
	)

	t.Run("Should be able to create a pack issuer NFT", func(t *testing.T) {

		// Assumes issuer is deployer of exampleNFT
		script = templates.GenerateCreatePackIssuerTx(pdsAddress)
		tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)

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
	script := templates.GenerateBorrowNFTScript(nftAddress, exampleNFTAddress)
	executeScriptAndCheck(
		t, b,
		script,
		[][]byte{
			jsoncdc.MustEncode(cadence.NewAddress(exampleNFTAddress)),
			jsoncdc.MustEncode(cadence.NewUInt64(0)),
		},
	)

	// Assumes issuer is deployer of exampleNFT
	// CreatePackIssuerTx
	script = templates.GenerateCreatePackIssuerTx(pdsAddress)
	tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)

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

	t.Run("Should be able to link NFT provider capability", func(t *testing.T) {

		// Assumes issuer is deployer of exampleNFT
		script = templates.GenerateLinkExampleNFTProviderCapTx(nftAddress, exampleNFTAddress)
		tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
		// Set argument: NFT provider path
		tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})

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

	t.Run("Should be able to set pack issuer capability", func(t *testing.T) {

		// Assumes issuer is deployer of exampleNFT
		script = templates.GenerateSetPackIssuerCapTx(pdsAddress)
		tx := createTxWithTemplateAndAuthorizer(b, script, pdsAddress)
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
		script = templates.GenerateCreateDistributionTx(pdsAddress, exampleNFTAddress, iPackNFTAddress, nftAddress)
		tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
		// Set argument: issuer address
		tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})
		tx.AddArgument(cadence.String("TestDistribution"))
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

		script = templates.GenerateGetDistTitleScript(pdsAddress)
		title := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.UInt64(1))})
		assert.Equal(t, cadence.String("TestDistribution"), title)
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
	script := templates.GenerateBorrowNFTScript(nftAddress, exampleNFTAddress)
	executeScriptAndCheck(
		t, b,
		script,
		[][]byte{
			jsoncdc.MustEncode(cadence.NewAddress(exampleNFTAddress)),
			jsoncdc.MustEncode(cadence.NewUInt64(0)),
		},
	)

	// Assumes issuer is deployer of exampleNFT
	// CreatePackIssuerTx
	script = templates.GenerateCreatePackIssuerTx(pdsAddress)
	tx := createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)

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

	// Assumes issuer is deployer of exampleNFT
	script = templates.GenerateLinkExampleNFTProviderCapTx(nftAddress, exampleNFTAddress)
	tx = createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
	// Set argument: NFT provider path
	tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})

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

	// Assumes issuer is deployer of exampleNFT
	script = templates.GenerateSetPackIssuerCapTx(pdsAddress)
	tx = createTxWithTemplateAndAuthorizer(b, script, pdsAddress)
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
	script = templates.GenerateCreateDistributionTx(pdsAddress, exampleNFTAddress, iPackNFTAddress, nftAddress)
	tx = createTxWithTemplateAndAuthorizer(b, script, exampleNFTAddress)
	// Set argument: issuer address
	tx.AddArgument(cadence.Path{Domain: "private", Identifier: "exampleNFTprovider"})
	tx.AddArgument(cadence.String("TestDistribution"))
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

	script = templates.GenerateGetDistTitleScript(pdsAddress)
	title := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.UInt64(1))})
	assert.Equal(t, cadence.String("TestDistribution"), title)

	t.Run("Should be able to mint a pack NFT", func(t *testing.T) {
		// Assumes issuer is deployer of exampleNFT
		script = templates.GenerateMintPackNFTTx(pdsAddress, exampleNFTAddress, nftAddress)
		tx := createTxWithTemplateAndAuthorizer(b, script, pdsAddress)
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

		script = templates.GeneratePackNFTTotalSupply(exampleNFTAddress)
		supply := executeScriptAndCheck(t, b, script, nil)
		assert.Equal(t, cadence.NewUInt64(1), supply)
	})
}
