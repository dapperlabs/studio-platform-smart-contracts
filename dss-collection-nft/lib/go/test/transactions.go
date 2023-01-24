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

func createCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	collectionGroupDescription string,
	typeName string,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(createCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.String(collectionGroupName))
	tx.AddArgument(cadence.String(collectionGroupDescription))
	tx.AddArgument(cadence.String(typeName))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)
	collectionGroupId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)

	return collectionGroupId
}

func createTimeBoundCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupName string,
	collectionGroupDescription string,
	typeName string,
	startTime int,
	endTime int,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(createTimeBoundCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.String(collectionGroupName))
	tx.AddArgument(cadence.String(collectionGroupDescription))
	tx.AddArgument(cadence.String(typeName))
	tx.AddArgument(cadence.UFix64(startTime))
	tx.AddArgument(cadence.UFix64(endTime))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)

	collectionGroupId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)

	return collectionGroupId
}

func closeCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	id uint64,
) {
	tx := flow.NewTransaction().
		SetScript(closeCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.UInt64(id))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)
}

func addNFTToCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupID uint64,
	nftID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(addNFTToCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.UInt64(collectionGroupID))
	tx.AddArgument(cadence.UInt64(nftID))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)
}

func mintNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	recipientAddress string,
	collectionGroupID uint64,
	completedBy string,
	level uint8,
) {
	tx := flow.NewTransaction().
		SetScript(mintDSSCollectionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.Address(flow.HexToAddress(recipientAddress)))
	tx.AddArgument(cadence.UInt64(collectionGroupID))
	tx.AddArgument(cadence.String(completedBy))
	tx.AddArgument(cadence.UInt8(level))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)
}
