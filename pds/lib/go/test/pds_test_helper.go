package test

import (
	"encoding/hex"
	studioPlatformContracts "github.com/dapperlabs/studio-platform-smart-contracts/lib/go/contracts"
	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Deploys the NonFungibleToken, MetadataViews, and ExampleNFT contracts to new accounts
// and returns their addresses
func deployPDSContracts(
	t *testing.T,
	b *emulator.Blockchain,
	exampleNFTAccountKey *flow.AccountKey,
	iPackNFTAccountKey *flow.AccountKey,
) (flow.Address, flow.Address, flow.Address, flow.Address, flow.Address) {

	// 1. Deploy utility contracts
	nftAddress := deploy(t, b, "NonFungibleToken", contracts.NonFungibleToken())
	metadataAddress := deploy(t, b, "MetadataViews", contracts.MetadataViews(flow.HexToAddress(emulatorFTAddress), nftAddress))
	exampleNFTAddress := deploy(
		t, b,
		"ExampleNFT",
		contracts.ExampleNFT(nftAddress, metadataAddress),
		exampleNFTAccountKey,
	)

	iPackNFTAddress := deploy(
		t, b,
		"IPackNFT",
		studioPlatformContracts.IPackNFT(nftAddress),
		iPackNFTAccountKey,
	)

	// 2. Deploy Pack NFT contract
	packNFTAddress := deployPackNftContract(t, b, nftAddress, iPackNFTAddress)

	// 3. Deploy PDS contract
	pdsAddress := deployPDSContract(t, b, nftAddress, iPackNFTAddress)

	return nftAddress, metadataAddress, exampleNFTAddress, packNFTAddress, pdsAddress
}

func deployPackNftContract(t *testing.T, b *emulator.Blockchain, nftAddress, iPackNFTAddress flow.Address) flow.Address {
	accountKeys := test.AccountKeyGenerator()

	// set up PackNFT account
	PackNftAccountKey, PackNftSigner := accountKeys.NewWithSigner()
	PackNftAddress, _ := b.CreateAccount([]*flow.AccountKey{PackNftAccountKey}, nil)
	PackNftCode := LoadPackNFT(nftAddress, iPackNFTAddress)
	fundAccount(t, b, PackNftAddress, defaultAccountFunding)

	packNFTencodedStr := hex.EncodeToString(PackNftCode)
	txBytes, _ := LoadPackNFTDeployScript()

	tx1 := createTxWithTemplateAndAuthorizer(b, txBytes, PackNftAddress)
	_ = tx1.AddArgument(cadence.String("PackNFT"))
	_ = tx1.AddArgument(cadence.String(packNFTencodedStr))
	_ = tx1.AddArgument(cadence.Path{Domain: "storage", Identifier: "PackNFTCollection"})
	_ = tx1.AddArgument(cadence.Path{Domain: "public", Identifier: "PackNFTCollectionPub"})
	_ = tx1.AddArgument(cadence.Path{Domain: "public", Identifier: "PackNFTIPackNFTCollectionPub"})
	_ = tx1.AddArgument(cadence.Path{Domain: "storage", Identifier: "PackNFTOperator"})
	_ = tx1.AddArgument(cadence.Path{Domain: "private", Identifier: "PackNFTOperatorPriv"})
	_ = tx1.AddArgument(cadence.String("0.1.0"))

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, PackNftAddress},
		[]crypto.Signer{signer, PackNftSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return PackNftAddress
}

func deployPDSContract(t *testing.T, b *emulator.Blockchain, nftAddress, iPackNFTAddress flow.Address) flow.Address {
	accountKeys := test.AccountKeyGenerator()

	// set up PackNFT account
	PDSAccountKey, PDSSigner := accountKeys.NewWithSigner()
	PDSAddress, _ := b.CreateAccount([]*flow.AccountKey{PDSAccountKey}, nil)
	PDSCode := LoadPDS(nftAddress, iPackNFTAddress)
	fundAccount(t, b, PDSAddress, defaultAccountFunding)

	PDSEncodedStr := hex.EncodeToString(PDSCode)
	script, _ := LoadPDSDeployScript()
	tx1 := createTxWithTemplateAndAuthorizer(b, script, PDSAddress)

	tx1.
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	_ = tx1.AddArgument(cadence.String("PDS"))
	_ = tx1.AddArgument(cadence.String(PDSEncodedStr))
	_ = tx1.AddArgument(cadence.Path{Domain: "storage", Identifier: "PDSPackIssuer"})
	_ = tx1.AddArgument(cadence.Path{Domain: "public", Identifier: "PDSPackIssuerCapRecv"})
	_ = tx1.AddArgument(cadence.Path{Domain: "storage", Identifier: "PDSDistCreator"})
	_ = tx1.AddArgument(cadence.Path{Domain: "private", Identifier: "PDSDistCap"})
	_ = tx1.AddArgument(cadence.Path{Domain: "storage", Identifier: "PDSDistManager"})
	_ = tx1.AddArgument(cadence.String("0.1.0"))

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, PDSAddress},
		[]crypto.Signer{signer, PDSSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return PDSAddress
}
