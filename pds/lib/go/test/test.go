package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/rs/zerolog"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-emulator/adapters"
	"github.com/onflow/flow-emulator/convert"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-emulator/types"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

const (
	emulatorFTAddress = "ee82856bf20e2aa6"
)

// Sets up testing and emulator objects and initialize the emulator default addresses
func newTestSetup(t *testing.T) (*emulator.Blockchain, *test.AccountKeys) {
	// Set for parallel processing
	t.Parallel()

	// Create a new emulator instance
	b := newBlockchain()

	// Create a new account key generator object to generate keys
	// for test accounts
	accountKeys := test.AccountKeyGenerator()

	return b, accountKeys
}

// newBlockchain returns an emulator blockchain for testing.
func newBlockchain(opts ...emulator.Option) *emulator.Blockchain {
	b, err := emulator.New(
		append(
			[]emulator.Option{
				emulator.WithStorageLimitEnabled(false),
			},
			opts...,
		)...,
	)
	if err != nil {
		panic(err)
	}
	return b
}

// Create a new, empty account for testing
// and return the address, public keys, and signer objects
func newAccountWithAddress(b *emulator.Blockchain, accountKeys *test.AccountKeys) (flow.Address, *flow.AccountKey, crypto.Signer) {
	newAccountKey, newSigner := accountKeys.NewWithSigner()
	logger := zerolog.Nop()
	adapter := adapters.NewSDKAdapter(&logger, b)
	newAddress, _ := adapter.CreateAccount(context.Background(), []*flow.AccountKey{newAccountKey}, nil)

	return newAddress, newAccountKey, newSigner
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

// Create a transaction object with the specified address as the authorizer
func createTxWithTemplateAndAuthorizer(
	b *emulator.Blockchain,
	script []byte,
	authorizerAddress flow.Address,
) *flow.Transaction {

	tx := flow.NewTransaction().
		SetScript(script).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(authorizerAddress)

	return tx
}

// signAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// before submitting it to the emulator.
//
// If the private keys do not match up with the addresses, the transaction will not succeed.
//
// The shouldRevert parameter indicates whether the transaction should fail or not.
//
// This function asserts the correct result and commits the block if it passed.
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

	Submit(t, b, tx, shouldRevert)
}

// signAndSubmitWithResult signs a transaction with an array of signers and adds their signatures to the transaction
// before submitting it to the emulator; it returns the transaction result.
//
// If the private keys do not match up with the addresses, the transaction will not succeed.
//
// The shouldRevert parameter indicates whether the transaction should fail or not.
//
// This function asserts the correct result and commits the block if it passed.
func signAndSubmitWithResult(
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

	return Submit(t, b, tx, shouldRevert)
}

// Submit submits a transaction and checks if it fails or not.
func Submit(
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

// Read a file from the specified path
func readFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// CadenceUFix64 returns a UFix64 value
func CadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)

	if err != nil {
		panic(err)
	}

	return newValue
}

func bytesToCadenceArray(b []byte) cadence.Array {
	values := make([]cadence.Value, len(b))

	for i, v := range b {
		values[i] = cadence.NewUInt8(v)
	}

	return cadence.NewArray(values)
}

// assertEqual asserts that two objects are equal.
//
//	assertEqual(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func assertEqual(t *testing.T, expected, actual interface{}) bool {

	if assert.ObjectsAreEqual(expected, actual) {
		return true
	}

	message := fmt.Sprintf(
		"Not equal: \nexpected: %s\nactual  : %s",
		expected,
		actual,
	)

	return assert.Fail(t, message)
}

func toJson(t *testing.T, target cadence.Value) string {
	actualJSONBytes, err := jsoncdc.Encode(target)
	require.NoError(t, err)
	actualJSON := string(actualJSONBytes)
	return actualJSON
}
