package test

import (
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	fttemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
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
	productName string,
	metadata map[string]string,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(createCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.String(collectionGroupName))
	tx.AddArgument(cadence.String(collectionGroupDescription))
	tx.AddArgument(cadence.String(productName))
	tx.AddArgument(metadataDict(metadata))

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
	productName string,
	endTime int,
	metadata map[string]string,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(createTimeBoundCollectionGroupTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.String(collectionGroupName))
	tx.AddArgument(cadence.String(collectionGroupDescription))
	tx.AddArgument(cadence.String(productName))
	tx.AddArgument(cadence.UFix64(endTime))
	tx.AddArgument(metadataDict(metadata))

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

func createSlot(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionGroupID uint64,
	logicalOperator string,
	required bool,
	metadata map[string]string,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(createSlotTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.UInt64(collectionGroupID))
	tx.AddArgument(cadence.String(logicalOperator))
	tx.AddArgument(cadence.Bool(required))
	tx.AddArgument(metadataDict(metadata))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)
	slotId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)

	return slotId
}

func createItemInSlot(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	itemID string,
	points uint64,
	itemType string,
	comparator string,
	slotID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(createItemInSlotTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.String(itemID))
	tx.AddArgument(cadence.UInt64(points))
	tx.AddArgument(cadence.String(itemType))
	tx.AddArgument(cadence.String(comparator))
	tx.AddArgument(cadence.UInt64(slotID))

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
	completionAddress string,
	level uint8,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(mintDSSCollectionTransaction(contracts)).
		SetGasLimit(999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.Address(flow.HexToAddress(recipientAddress)))
	tx.AddArgument(cadence.UInt64(collectionGroupID))
	tx.AddArgument(cadence.String(completionAddress))
	tx.AddArgument(cadence.UInt8(level))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)

	nftId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)
	return nftId
}

func mintNFTAndRecordCompletedWith(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	recipientAddress string,
	collectionGroupID uint64,
	completionAddress string,
	level uint8,
	nftIDs []uint64,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(mintDSSCollectionAndRecordTransaction(contracts)).
		SetGasLimit(999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.Address(flow.HexToAddress(recipientAddress)))
	tx.AddArgument(cadence.UInt64(collectionGroupID))
	tx.AddArgument(cadence.String(completionAddress))
	tx.AddArgument(cadence.UInt8(level))

	cadenceNftIDs := make([]cadence.Value, len(nftIDs))
	for i, nftID := range nftIDs {
		cadenceNftIDs[i] = cadence.UInt64(nftID)
	}

	tx.AddArgument(cadence.NewArray(cadenceNftIDs))
	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)

	nftId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)
	return nftId
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
		AddAuthorizer(contracts.DSSCollectionAddress)
	tx.AddArgument(cadence.Address(flow.HexToAddress(recipientAddress)))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)

	nftId := txResult.Events[0].Value.Fields[0].ToGoValue().(uint64)
	return nftId
}

func completedCollectionGroup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	collectionID uint64,
	userAddress string,
	nftIDs []uint64,
) {
	tx := flow.NewTransaction().
		SetScript(setCompletedCollectionGroup(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DSSCollectionAddress)

	tx.AddArgument(cadence.UInt64(collectionID))
	tx.AddArgument(cadence.Address(flow.HexToAddress(userAddress)))

	cadenceNftIDs := make([]cadence.Value, len(nftIDs))
	for i, nftID := range nftIDs {
		cadenceNftIDs[i] = cadence.UInt64(nftID)
	}
	tx.AddArgument(cadence.NewArray(cadenceNftIDs))
	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DSSCollectionAddress},
		[]crypto.Signer{signer, contracts.DSSCollectionSigner},
		shouldRevert,
	)
}

func transferNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	toAddress string,
	nftID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(transferNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)
	tx.AddArgument(cadence.UInt64(nftID))
	tx.AddArgument(cadence.Address(flow.HexToAddress(toAddress)))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{signer, userSigner},
		shouldRevert,
	)
}
