package test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/adapters"
	"github.com/onflow/flow-emulator/convert"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-emulator/types"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/rs/zerolog"
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
	NFTAddress                   flow.Address
	ExampleNFTAddress            flow.Address
	BurnerAddress                flow.Address
	ViewResolverAddress          flow.Address
	MetadataViewAddress          flow.Address
	NFTProviderAggregatorAddress flow.Address
	NFTProviderAggregatorSigner  crypto.Signer
}

func deployNFTContracts(t *testing.T, b *emulator.Blockchain) (flow.Address, flow.Address, flow.Address, flow.Address) {
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	resolverAddress := deploy(t, adapter, "ViewResolver", nftcontracts.ViewResolver(), b.ServiceKey().AccountKey())

	viewResolverRe := regexp.MustCompile(viewResolverAddressPlaceholder)
	nftCode := viewResolverRe.ReplaceAll(readFile(NonFungibleTokenPath), []byte("0x"+resolverAddress.String()))
	nftAddress, err := adapter.CreateAccount(context.Background(),
		nil,
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

	metadataViewsAddr := deploy(t, adapter, "MetadataViews", nftcontracts.MetadataViews(ftAddress.String(), nftAddress.String(), resolverAddress.String()), b.ServiceKey().AccountKey())

	burnerAddr := deploy(t, adapter, "Burner", readFile(BurnerContractPath), b.ServiceKey().AccountKey())

	return nftAddress, metadataViewsAddr, resolverAddress, burnerAddr
}

func NFTProviderAggregatorDeployContracts(t *testing.T, b *emulator.Blockchain) Contracts {
	accountKeys := test.AccountKeyGenerator()

	nftAddress, metadataViewsAddr, resolverAddress, burnerAddr := deployNFTContracts(t, b)
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)

	nftProviderAggregatorAccountKey, nftProviderAggregatorSigner := accountKeys.NewWithSigner()
	nftProviderAggregator := LoadNftProviderAggregator(nftAddress, metadataViewsAddr, resolverAddress, burnerAddr)

	nftProviderAggregatorAddress, err := adapter.CreateAccount(context.Background(),
		[]*flow.AccountKey{nftProviderAggregatorAccountKey},
		nil,
	)
	require.NoError(t, err)

	addContractScript := `
	transaction(name: String, code: String) {
		prepare(signer: auth(AddContract) &Account) {
			signer.contracts.add(name: name, code: code.decodeHex())
		}
	}`

	// Deploy NFTProviderAggregator contract
	nftProviderAggregatorContract := sdktemplates.Contract{
		Name:   "NFTProviderAggregator",
		Source: string(nftProviderAggregator),
	}
	tx := flow.NewTransaction().
		AddRawArgument(jsoncdc.MustEncode(cadence.String(nftProviderAggregatorContract.Name))).
		AddRawArgument(jsoncdc.MustEncode(cadence.String(nftProviderAggregatorContract.SourceHex()))).
		AddAuthorizer(nftProviderAggregatorAddress)
	tx.SetScript([]byte(addContractScript))
	tx.SetComputeLimit(200).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signer, err := b.ServiceKey().Signer()
	assert.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, nftProviderAggregatorAddress},
		[]crypto.Signer{signer, nftProviderAggregatorSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	// Deploy ExampleNFT contract
	exampleNFTcontract := sdktemplates.Contract{
		Name:   "ExampleNFT",
		Source: string(nftcontracts.ExampleNFT(nftAddress, metadataViewsAddr, resolverAddress)),
	}
	tx = flow.NewTransaction().
		AddRawArgument(jsoncdc.MustEncode(cadence.String(exampleNFTcontract.Name))).
		AddRawArgument(jsoncdc.MustEncode(cadence.String(exampleNFTcontract.SourceHex()))).
		AddAuthorizer(nftProviderAggregatorAddress)
	tx.SetScript([]byte(addContractScript))
	tx.SetComputeLimit(200).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, nftProviderAggregatorAddress},
		[]crypto.Signer{signer, nftProviderAggregatorSigner},
		false,
	)

	_, err = b.CommitBlock()
	require.NoError(t, err)

	return Contracts{
		nftAddress,
		nftProviderAggregatorAddress,
		burnerAddr,
		resolverAddress,
		metadataViewsAddr,
		nftProviderAggregatorAddress,
		nftProviderAggregatorSigner,
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
	contents, err := os.ReadFile(path)
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

func setupExampleNft(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	contracts Contracts,
) {
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, ExampleNftSetupPath)).
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

func metadataDict(dict map[string]interface{}) cadence.Dictionary {
	pairs := []cadence.KeyValuePair{}

	for key, value := range dict {
		cadenceKey, _ := cadence.NewString(key)
		// if value is a string, we need to convert it to a cadence string
		// if value is an array of strings, we need to convert it to a cadence array of cadence strings
		if strValue, ok := value.(string); ok {
			cadenceValue, _ := cadence.NewString(strValue)
			pairs = append(pairs, cadence.KeyValuePair{Key: cadenceKey, Value: cadenceValue})
		} else if strArrayValue, ok := value.([]string); ok {
			cadenceValueArray := make([]cadence.Value, len(strArrayValue))
			for index, strValue := range strArrayValue {
				cadenceValueArray[index], _ = cadence.NewString(strValue)
			}
			cadenceValue := cadence.NewArray(cadenceValueArray)
			pairs = append(pairs, cadence.KeyValuePair{Key: cadenceKey, Value: cadenceValue})
		}
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
