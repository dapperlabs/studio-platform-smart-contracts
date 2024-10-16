package test

import (
	"context"
	"io/ioutil"
	"testing"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/adapters"
	"github.com/onflow/flow-emulator/convert"
	"github.com/onflow/flow-emulator/types"
	"github.com/rs/zerolog"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	flowTokenName         = "FlowToken"
	nonFungibleTokenName  = "NonFungibleToken"
	defaultAccountFunding = "1000.0"
)

var (
	ftAddress = flow.HexToAddress("ee82856bf20e2aa6")
)

type Contracts struct {
	NFTAddress           flow.Address
	MetadataViewsAddress flow.Address
	NFTLockerAddress     flow.Address
	NFTLockerSigner      crypto.Signer
	ExampleNFTAddress    flow.Address
}

func deployNFTContract(t *testing.T, b *emulator.Blockchain) (flow.Address, flow.Address) {
	resolverAddress := deploy(t, b, "ViewResolver", nftcontracts.ViewResolver(), b.ServiceKey().AccountKey())

	nftCode := nftcontracts.NonFungibleToken(resolverAddress.String())
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	nftAddress, err := adapter.CreateAccount(context.Background(), nil,
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

	return resolverAddress, nftAddress
}

func NFTLockerDeployContracts(t *testing.T, b *emulator.Blockchain) Contracts {
	accountKeys := test.AccountKeyGenerator()
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	resolverAddress, nftAddress := deployNFTContract(t, b)
	metadataViewsAddr := deploy(t, b, "MetadataViews", nftcontracts.MetadataViews(ftAddress.String(), nftAddress.String(), resolverAddress.String()))

	NFTLockerAccountKey, NFTLockerSigner := accountKeys.NewWithSigner()
	NFTLockerCode := LoadNFTLockerContract(nftAddress, metadataViewsAddr)

	ExampleNFTCode := nftcontracts.ExampleNFT(nftAddress, metadataViewsAddr, resolverAddress)

	NFTLockerAddress, err := adapter.CreateAccount(
		context.Background(),
		[]*flow.AccountKey{NFTLockerAccountKey},
		nil,
	)
	require.NoError(t, err)

	EscrowCode := LoadEscrowContract(nftAddress, metadataViewsAddr, NFTLockerAddress)

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	tx1 := sdktemplates.AddAccountContract(
		NFTLockerAddress,
		sdktemplates.Contract{
			Name:   "NFTLocker",
			Source: string(NFTLockerCode),
		},
	)

	tx1.
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, NFTLockerAddress},
		[]crypto.Signer{signer, NFTLockerSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	tx2 := sdktemplates.AddAccountContract(
		NFTLockerAddress,
		sdktemplates.Contract{
			Name:   "ExampleNFT",
			Source: string(ExampleNFTCode),
		},
	)

	tx2.
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signAndSubmit(
		t, b, tx2,
		[]flow.Address{b.ServiceKey().Address, NFTLockerAddress},
		[]crypto.Signer{signer, NFTLockerSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	tx3 := sdktemplates.AddAccountContract(
		NFTLockerAddress,
		sdktemplates.Contract{
			Name:   "Escrow",
			Source: string(EscrowCode),
		},
	)

	tx3.
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signAndSubmit(
		t, b, tx3,
		[]flow.Address{b.ServiceKey().Address, NFTLockerAddress},
		[]crypto.Signer{signer, NFTLockerSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return Contracts{
		nftAddress,
		metadataViewsAddr,
		NFTLockerAddress,
		NFTLockerSigner,
		NFTLockerAddress,
	}
}

// newEmulator returns a emulator object for testing
func newEmulator() *emulator.Blockchain {
	b, err := emulator.New(emulator.WithStorageLimitEnabled(false))
	if err != nil {
		panic(err)
	}
	return b
}

// Deploy a contract to a new account with the specified name, code, and keys
func deploy(
	t *testing.T,
	b *emulator.Blockchain,
	name string,
	code []byte,
	keys ...*flow.AccountKey,
) flow.Address {
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	address, err := adapter.CreateAccount(context.Background(),
		keys,
		[]sdktemplates.Contract{
			{
				Name:   name,
				Source: string(code),
			},
		},
	)
	assert.NoError(t, err)

	return address
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
) *types.TransactionResult {
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

	result := submit(t, b, tx, shouldRevert)
	return result
}

// submit submits a transaction and checks
// if it fails or not
func submit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	shouldRevert bool,
) *types.TransactionResult {
	// submit the signed transaction
	flowTx := convert.SDKTransactionToFlow(*tx)
	err := b.AddTransaction(*flowTx)
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
	return result
}

// executeScriptAndCheck executes a script and checks to make sure that it succeeded.
func executeScriptAndCheck(t *testing.T, b *emulator.Blockchain, script []byte, arguments [][]byte) cadence.Value {
	result, err := b.ExecuteScript(script, arguments)
	require.NoError(t, err)
	if !assert.True(t, result.Succeeded()) {
		t.Log(result.Error.Error())
	}
	return result.Value
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

// Simple error-handling wrapper for Flow account creation.
func createAccount(t *testing.T, b *emulator.Blockchain) (flow.Address, crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()
	accountKey, signer := accountKeys.NewWithSigner()

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	address, err := adapter.CreateAccount(context.Background(), []*flow.AccountKey{accountKey}, nil)
	require.NoError(t, err)

	return address, signer
}

func setupNFTLockerAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	contracts Contracts,
) {
	tx := flow.NewTransaction().
		SetScript(setupAccountTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{signer, userSigner},
		false,
	)
}

func setupExampleNFT(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	contracts Contracts,
) {
	tx := flow.NewTransaction().
		SetScript(setupExampleNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{signer, userSigner},
		false,
	)
}

func getLockedTokenData(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	id uint64,
) LockedData {
	script := readLockedTokenByIDScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.UInt64(id))})

	return parseLockedData(result)
}

func getInventoryData(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	address string,
) []uint64 {
	script := readInventoryScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.NewAddress(flow.HexToAddress(address)))})

	return parseInventoryData(result)
}
