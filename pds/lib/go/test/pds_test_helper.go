package test

import (
	"context"
	"encoding/hex"

	studioPlatformContracts "github.com/dapperlabs/studio-platform-smart-contracts/lib/go/contracts"
	studioPlatformTemplates "github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates"
	"github.com/onflow/flow-nft/lib/go/templates"
	"github.com/rs/zerolog"

	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/flow-emulator/adapters"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"
	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Deploys the NonFungibleToken, MetadataViews, and ExampleNFT contracts to new accounts
// and returns their addresses
func deployPDSContracts(
	t *testing.T,
	b *emulator.Blockchain,
	exampleNFTAccountKey *flow.AccountKey,
	exampleNFTSigner crypto.Signer,
	iPackNFTAccountKey *flow.AccountKey,
	pdsAccountKey *flow.AccountKey,
	pdsSigner crypto.Signer,
) (flow.Address, flow.Address, flow.Address, flow.Address, flow.Address) {

	// 1. Deploy utility contracts
	resolverAddress := deploy(t, b, "ViewResolver", nftcontracts.ViewResolver(), b.ServiceKey().AccountKey())

	nftAddress := deploy(t, b, "NonFungibleToken", nftcontracts.NonFungibleToken(resolverAddress.String()))

	metadataAddress := deploy(t, b, "MetadataViews", nftcontracts.MetadataViews(emulatorFTAddress, nftAddress.String(), resolverAddress.String()))
	exampleNFTAddress := deploy(
		t, b,
		"ExampleNFT",
		nftcontracts.ExampleNFT(nftAddress, metadataAddress, resolverAddress),
		exampleNFTAccountKey,
	)

	iPackNFTAddress := deploy(
		t, b,
		"IPackNFT",
		studioPlatformContracts.IPackNFT(nftAddress),
		iPackNFTAccountKey,
	)

	// 2. Deploy Pack NFT contract
	deployPackNftContract(t, b, nftAddress, iPackNFTAddress, exampleNFTAddress, metadataAddress, exampleNFTSigner)

	// 3. Deploy AllDay Pack NFT contract
	deployAllDayPackNftContract(t, b, nftAddress, ftAddress, iPackNFTAddress, metadataAddress)

	// 4. Deploy PDS contract
	pdsAddress := deployPDSContract(t, b, nftAddress, iPackNFTAddress, pdsAccountKey, pdsSigner)

	return nftAddress, metadataAddress, exampleNFTAddress, iPackNFTAddress, pdsAddress
}

func deployPackNftContract(t *testing.T, b *emulator.Blockchain, nftAddress, iPackNFTAddress, exampleNFTAddress, metadataAddress flow.Address, exampleNFTSigner crypto.Signer) {

	PackNftCode := studioPlatformContracts.PackNFT(nftAddress, iPackNFTAddress, metadataAddress)

	packNFTencodedStr := hex.EncodeToString(PackNftCode)
	txBytes := studioPlatformTemplates.GenerateDeployPackNFTTx(nftAddress, iPackNFTAddress)

	tx1 := createTxWithTemplateAndAuthorizer(b, txBytes, exampleNFTAddress)
	_ = tx1.AddArgument(cadence.String("PackNFT"))
	_ = tx1.AddArgument(cadence.String(packNFTencodedStr))
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PackNFTCollection"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainPublic, Identifier: "PackNFTCollectionPub"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainPublic, Identifier: "PackNFTIPackNFTCollectionPub"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PackNFTOperator"})
	// _ = tx1.AddArgument(cadence.Path{Domain: "private", Identifier: "PackNFTOperatorPriv"})
	_ = tx1.AddArgument(cadence.String("0.1.0"))

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, exampleNFTAddress},
		[]crypto.Signer{signer, exampleNFTSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)
}

func deployAllDayPackNftContract(t *testing.T, b *emulator.Blockchain, nftAddress, ftAddress, iPackNFTAddress, metaDataViewAddress flow.Address) flow.Address {
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	accountKeys := test.AccountKeyGenerator()

	// set up PackNFT account
	AllDayPackNftAccountKey, AllDayPackNftSigner := accountKeys.NewWithSigner()
	AllDayPackNftAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{AllDayPackNftAccountKey}, nil)
	PackNftCode := studioPlatformContracts.AllDayPackNFT(nftAddress, ftAddress, iPackNFTAddress, metaDataViewAddress, AllDayPackNftAddress)

	packNFTencodedStr := hex.EncodeToString(PackNftCode)
	txBytes := studioPlatformTemplates.GenerateDeployPackNFTTx(nftAddress, iPackNFTAddress)

	tx1 := createTxWithTemplateAndAuthorizer(b, txBytes, AllDayPackNftAddress)
	_ = tx1.AddArgument(cadence.String("PackNFT"))
	_ = tx1.AddArgument(cadence.String(packNFTencodedStr))
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PackNFTCollection"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainPublic, Identifier: "PackNFTCollectionPub"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainPublic, Identifier: "PackNFTIPackNFTCollectionPub"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PackNFTOperator"})
	// _ = tx1.AddArgument(cadence.Path{Domain: "private", Identifier: "PackNFTOperatorPriv"})
	_ = tx1.AddArgument(cadence.String("0.1.0"))

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, AllDayPackNftAddress},
		[]crypto.Signer{signer, AllDayPackNftSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return AllDayPackNftAddress
}

func deployPDSContract(
	t *testing.T,
	b *emulator.Blockchain,
	nftAddress, iPackNFTAddress flow.Address,
	pdsAccountKey *flow.AccountKey,
	pdsSigner crypto.Signer) flow.Address {
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	//accountKeys := test.AccountKeyGenerator()

	// set up PackNFT account
	//PDSAccountKey, PDSSigner := accountKeys.NewWithSigner()
	PDSAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{pdsAccountKey}, nil)
	PDSCode := studioPlatformContracts.PDS(nftAddress, iPackNFTAddress)

	PDSEncodedStr := hex.EncodeToString(PDSCode)
	script := studioPlatformTemplates.GenerateDeployPDSTx(nftAddress, iPackNFTAddress)
	tx1 := createTxWithTemplateAndAuthorizer(b, script, PDSAddress)

	tx1.
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	_ = tx1.AddArgument(cadence.String("PDS"))
	_ = tx1.AddArgument(cadence.String(PDSEncodedStr))
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PDSPackIssuer"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainPublic, Identifier: "PDSPackIssuerCapRecv"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PDSDistCreator"})
	// _ = tx1.AddArgument(cadence.Path{Domain: "private", Identifier: "PDSDistCap"})
	_ = tx1.AddArgument(cadence.Path{Domain: common.PathDomainStorage, Identifier: "PDSDistManager"})
	_ = tx1.AddArgument(cadence.String("0.1.0"))

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, PDSAddress},
		[]crypto.Signer{signer, pdsSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return PDSAddress
}

// Asserts that the ExampleNFT collection in the specified user's account
// is the expected length
func assertCollectionLength(
	t *testing.T,
	b *emulator.Blockchain,
	nftAddress, exampleNFTAddress, metadataAddress, collectionAddress flow.Address,
	expectedLength int,
) {
	script := templates.GenerateGetCollectionLengthScript(nftAddress, exampleNFTAddress, metadataAddress)
	actualLength := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.NewAddress(collectionAddress))})
	assert.Equal(t, cadence.NewInt(expectedLength), actualLength)
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

// Sets up an account with the generic royalty receiver in place of their Flow token receiver
func setupRoyaltyReceiver(
	t *testing.T,
	b *emulator.Blockchain,
	metadataAddress flow.Address,
	authorizerAddress flow.Address,
	authorizerSigner crypto.Signer,
) {

	script := templates.GenerateSetupAccountToReceiveRoyaltyScript(metadataAddress, flow.HexToAddress(emulatorFTAddress))
	tx := createTxWithTemplateAndAuthorizer(b, script, authorizerAddress)

	vaultPath := cadence.Path{Domain: common.PathDomainStorage, Identifier: "flowTokenVault"}
	tx.AddArgument(vaultPath)

	serviceSigner, _ := b.ServiceKey().Signer()

	signAndSubmit(
		t, b, tx,
		[]flow.Address{
			b.ServiceKey().Address,
			authorizerAddress,
		},
		[]crypto.Signer{
			serviceSigner,
			authorizerSigner,
		},
		false,
	)
}
