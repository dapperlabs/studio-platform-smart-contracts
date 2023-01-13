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

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address},
		[]crypto.Signer{b.ServiceKey().Signer()},
		false,
	)
}

func createCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	productPathDomain string,
	productPathIdentifier string,
) {
	tx := flow.NewTransaction().
		SetScript(createCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.String(collectionGroupName))
	tx.AddArgument(cadence.Path{
		Domain:     productPathDomain,
		Identifier: productPathIdentifier,
	})

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.DSSCollectionSigner},
		shouldRevert,
	)
}
