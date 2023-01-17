package test

import (
	"encoding/hex"
	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// Deploys the NonFungibleToken, MetadataViews, and ExampleNFT contracts to new accounts
// and returns their addresses
func deployNFTContracts(
	t *testing.T,
	b *emulator.Blockchain,
	exampleNFTAccountKey *flow.AccountKey,
) (flow.Address, flow.Address, flow.Address) {

	nftAddress := deploy(t, b, "NonFungibleToken", contracts.NonFungibleToken())
	metadataAddress := deploy(t, b, "MetadataViews", contracts.MetadataViews(flow.HexToAddress(emulatorFTAddress), nftAddress))

	exampleNFTAddress := deploy(
		t, b,
		"ExampleNFT",
		contracts.ExampleNFT(nftAddress, metadataAddress),
		exampleNFTAccountKey,
	)

	return nftAddress, metadataAddress, exampleNFTAddress
}

func deployPackNftContracts(t *testing.T, b *emulator.Blockchain) PackNftContracts {
	accountKeys := test.AccountKeyGenerator()

	nftAddress := deployNFTContract(t, b)
	iPackCode := LoadIPackNFT(nftAddress)
	iPackAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "IPackNFT",
			Source: string(iPackCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// set up PackNFT account
	PackNftAccountKey, PackNftSigner := accountKeys.NewWithSigner()
	PackNftAddress, _ := b.CreateAccount([]*flow.AccountKey{PackNftAccountKey}, nil)
	PackNftCode := LoadPackNFT(nftAddress, iPackAddress)
	fundAccount(t, b, PackNftAddress, defaultAccountFunding)
	require.NoError(t, err)

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

	return PackNftContracts{
		nftAddress,
		PackNftAddress,
		PackNftSigner,
	}
}

func deployPDSContracts(t *testing.T, b *emulator.Blockchain) PdsContracts {
	accountKeys := test.AccountKeyGenerator()

	nftAddress := deployNFTContract(t, b)
	iPackCode := LoadIPackNFT(nftAddress)
	iPackAddress, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "IPackNFT",
			Source: string(iPackCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// set up PackNFT account
	PDSAccountKey, PDSSigner := accountKeys.NewWithSigner()
	PDSAddress, _ := b.CreateAccount([]*flow.AccountKey{PDSAccountKey}, nil)
	PDSCode := LoadPDS(nftAddress, iPackAddress)
	fundAccount(t, b, PDSAddress, defaultAccountFunding)
	require.NoError(t, err)

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

	return PdsContracts{
		nftAddress,
		PDSAddress,
		PDSSigner,
	}
}
