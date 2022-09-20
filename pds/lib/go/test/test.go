package test

import (
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/onflow/flow-go-sdk"
)

const (
	flowTokenName         = "FlowToken"
	nonFungibleTokenName  = "NonFungibleToken"
	defaultAccountFunding = "1000.0"
)

var (
	ftAddress        = flow.HexToAddress("ee82856bf20e2aa6")
	flowTokenAddress = flow.HexToAddress("0ae53cb6e3f42a79")
)

type PackNftContracts struct {
	NFTAddress     flow.Address
	PackNftAddress flow.Address
	PackNftSigner  crypto.Signer
}

type PdsContracts struct {
	NFTAddress flow.Address
	PdsAddress flow.Address
	PdsSigner  crypto.Signer
}

func deployNFTContract(t *testing.T, b *emulator.Blockchain) flow.Address {
	nftCode := nftcontracts.NonFungibleToken()
	nftAddress, err := b.CreateAccount(nil,
		[]sdktemplates.Contract{
			{
				Name:   nonFungibleTokenName,
				Source: string(nftCode),
			},
		},
	)
	require.NoError(t, err)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return nftAddress
}

func PackNftDeployContracts(t *testing.T, b *emulator.Blockchain) PackNftContracts {
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
	latestBlock, _ := b.GetLatestBlock()
	blockId := latestBlock.ID()

	tx1 := flow.NewTransaction().
		SetScript(txBytes).SetGasLimit(100).
		SetProposalKey(PackNftAddress, PackNftAccountKey.Index, PackNftAccountKey.SequenceNumber).
		SetReferenceBlockID(sdk.Identifier(blockId)).
		SetPayer(PackNftAddress).
		AddAuthorizer(PackNftAddress)

	tx1.
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

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

func PDSDeployContracts(t *testing.T, b *emulator.Blockchain) PdsContracts {
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
	txBytes, _ := LoadPDSDeployScript()
	latestBlock, _ := b.GetLatestBlock()
	blockId := latestBlock.ID()

	tx1 := flow.NewTransaction().
		SetScript(txBytes).SetGasLimit(100).
		SetProposalKey(PDSAddress, PDSAccountKey.Index, PDSAccountKey.SequenceNumber).
		SetReferenceBlockID(sdk.Identifier(blockId)).
		SetPayer(PDSAddress).
		AddAuthorizer(PDSAddress)

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

// newEmulator returns a emulator object for testing
func newEmulator() *emulator.Blockchain {
	b, err := emulator.NewBlockchain()
	if err != nil {
		panic(err)
	}
	return b
}

// signAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// Then submits the transaction to the emulator. If the private keys don't match up with the addresses,
// the transaction will not succeed.
// shouldRevert parameter indicates whether the transaction should fail or not
// This function asserts the correct result and commits the block if it passed
func signAndSubmit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) {
	// sign transaction with each signer
	for i := len(signerAddresses) - 1; i >= 0; i-- {
		signerAddress := signerAddresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.SignEnvelope(signerAddress, 0, signer)
			assert.NoError(t, err)
		} else {
			err := tx.SignPayload(signerAddress, 0, signer)
			assert.NoError(t, err)
		}
	}

	submit(t, b, tx, shouldRevert)
}

// submit submits a transaction and checks
// if it fails or not
func submit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	shouldRevert bool,
) {
	// submit the signed transaction
	err := b.AddTransaction(*tx)
	require.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	require.NoError(t, err)

	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	}

	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

// readFile reads a file from the file system
// and returns its contents
func readFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// cadenceUFix64 returns a UFix64 value
func cadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}

	return newValue
}
