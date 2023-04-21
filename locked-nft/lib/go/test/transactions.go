package test

import (
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	fttemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------
func fundAccount(
	t *testing.T,
	b *emulator.Blockchain,
	receiverAddress flow.Address,
	amount string,
) {
	script := fttemplates.GenerateMintTokensScript(
		ftAddress,
		flowTokenAddress,
		flowTokenName,
	)

	tx := flow.NewTransaction().
		SetScript(script).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(b.ServiceKey().Address)

	tx.AddArgument(cadence.NewAddress(receiverAddress))
	tx.AddArgument(cadenceUFix64(amount))

	signer, _ := b.ServiceKey().Signer()

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address},
		[]crypto.Signer{signer},
		false,
	)
}

func mintExampleNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	recipientAddress string,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(mintExampleNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.LockedNFTAddress)
	tx.AddArgument(cadence.Address(flow.HexToAddress(recipientAddress)))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.LockedNFTAddress},
		[]crypto.Signer{signer, contracts.LockedNFTSigner},
		shouldRevert,
	)

	nftId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)
	return nftId
}
