package test

import (
	"context"
	"github.com/onflow/flow-emulator/convert"
	"github.com/rs/zerolog"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-emulator/adapters"
	emulator "github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-emulator/types"
	ftcontracts "github.com/onflow/flow-ft/lib/go/contracts"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
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
	fungibleTokenName     = "FungibleToken"
	defaultAccountFunding = "1000.0"
)

var (
	ftAddress        = flow.HexToAddress("ee82856bf20e2aa6")
	flowTokenAddress = flow.HexToAddress("0ae53cb6e3f42a79")
)

type Contracts struct {
	BurnerAddress       flow.Address
	ViewResolverAddress flow.Address
	NFTAddress          flow.Address
	MetadataViewAddress flow.Address
	GolazosAddress      flow.Address
	GolazosSigner       crypto.Signer
	FtAddress           flow.Address
}

func GolazosDeployContracts(t *testing.T, b *emulator.Blockchain) Contracts {
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	resolverAddress := deploy(t, adapter, "ViewResolver", nftcontracts.ViewResolver(), b.ServiceKey().AccountKey())
	burnerAddress := deploy(t, adapter, "Burner", readFile(BurnerPath), b.ServiceKey().AccountKey())

	nftCode := nftcontracts.NonFungibleToken(resolverAddress.String())
	nftAddress := deploy(t, adapter, nonFungibleTokenName, nftCode, b.ServiceKey().AccountKey())

	ftCode := ftcontracts.FungibleToken(resolverAddress.String(), burnerAddress.String())
	ftAddress := deploy(t, adapter, fungibleTokenName, ftCode, b.ServiceKey().AccountKey())

	accountKeys := test.AccountKeyGenerator()

	metadataCode := nftcontracts.MetadataViews(ftAddress.String(), nftAddress.String(), resolverAddress.String())
	metadataViewsAddr := deploy(t, adapter, "MetadataViews", metadataCode, b.ServiceKey().AccountKey())

	GolazosAccountKey, GolazosSigner := accountKeys.NewWithSigner()
	golazosAddress, err := adapter.CreateAccount(context.Background(),
		[]*sdk.AccountKey{GolazosAccountKey},
		nil,
	)
	require.NoError(t, err)

	GolazosCode := LoadGolazos(nftAddress, metadataViewsAddr, ftAddress, resolverAddress)
	GolazosCode = []byte(strings.ReplaceAll(string(GolazosCode), royaltyAddressPlaceholder, "0x"+golazosAddress.String()))

	//fundAccount(t, b, GolazosAddress, defaultAccountFunding)

	tx1 := sdktemplates.AddAccountContract(
		golazosAddress,
		sdktemplates.Contract{
			Name:   "Golazos",
			Source: string(GolazosCode),
		},
	)

	tx1.
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx1,
		[]flow.Address{b.ServiceKey().Address, golazosAddress},
		[]crypto.Signer{signer, GolazosSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return Contracts{
		burnerAddress,
		resolverAddress,
		nftAddress,
		metadataViewsAddr,
		golazosAddress,
		GolazosSigner,
		ftAddress,
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

	return submit(t, b, tx, shouldRevert)
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

// cadenceUFix64 returns a UFix64 value
func cadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}

	return newValue
}

// Simple error-handling wrapper for Flow account creation.
func createAccount(t *testing.T, b *emulator.Blockchain) (sdk.Address, crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()
	accountKey, signer := accountKeys.NewWithSigner()

	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	address, err := adapter.CreateAccount(context.Background(), []*sdk.AccountKey{accountKey}, nil)
	require.NoError(t, err)

	return address, signer
}

func setupGolazos(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	contracts Contracts,
) {
	tx := flow.NewTransaction().
		SetScript(loadGolazosSetupAccountTransaction(contracts)).
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

func setupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	address flow.Address,
	signer crypto.Signer,
	contracts Contracts,
) (sdk.Address, crypto.Signer) {
	setupGolazos(t, b, address, signer, contracts)
	//fundAccount(t, b, address, defaultAccountFunding)

	return address, signer
}

func metadataDict(dict map[string]string) cadence.Dictionary {
	pairs := []cadence.KeyValuePair{}

	for key, value := range dict {
		cadenceKey, _ := cadence.NewString(key)
		cadenceValue, _ := cadence.NewString(value)
		pairs = append(pairs, cadence.KeyValuePair{Key: cadenceKey, Value: cadenceValue})
	}

	return cadence.NewDictionary(pairs)
}

// Deploy a contract to a new account with the specified name, code, and keys
func deploy(
	t *testing.T,
	adapter *adapters.SDKAdapter,
	name string,
	code []byte,
	keys ...*flow.AccountKey,
) flow.Address {
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
