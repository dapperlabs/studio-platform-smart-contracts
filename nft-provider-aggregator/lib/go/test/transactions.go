package test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-emulator/types"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/require"
)

func mintExampleNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipientAddress flow.Address,
	shouldRevert bool,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, ExampleNftMintPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	if err := tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes())); err != nil {
		t.Error(err)
	}

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	result := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
	nftID := uint64(0)
	for _, event := range result.Events {
		if strings.Contains(event.Type, "Deposited") {
			if v := cadence.FieldsMappedByName(event.Value)["id"]; v != nil {
				nftID = GetFieldValue(v).(uint64)
			}
		}
	}
	return nftID
}

func bootstrapAggregatorResource(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	nftIdentifier string,
	supplierAddresses []flow.Address,
	capabilityPublicationIDs []string,
	shouldRevert bool,
) {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, BootstrapAggregatorPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Add arguments
	cadenceNftIdentifier, _ := cadence.NewString(nftIdentifier)
	if err := tx.AddArgument(cadenceNftIdentifier); err != nil {
		t.Error(err)
	}
	var cadenceSuppliers []cadence.Value
	for _, supplierAddress := range supplierAddresses {
		cadenceSuppliers = append(cadenceSuppliers, cadence.BytesToAddress(supplierAddress.Bytes()))
	}
	if err := tx.AddArgument(cadence.NewArray(cadenceSuppliers)); err != nil {
		t.Error(err)
	}
	var cadenceCapabilityPublicationIDs []cadence.Value
	for _, capabilityPublicationID := range capabilityPublicationIDs {
		cadenceCapabilityPublicationID, _ := cadence.NewString(capabilityPublicationID)
		cadenceCapabilityPublicationIDs = append(cadenceCapabilityPublicationIDs, cadenceCapabilityPublicationID)
	}
	if err := tx.AddArgument(cadence.NewArray(cadenceCapabilityPublicationIDs)); err != nil {
		t.Error(err)
	}

	// Sign and submit
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
}

func bootstrapSupplierResource(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	capabilityPublicationID string,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, BootstrapSupplierPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.BytesToAddress(contracts.NFTProviderAggregatorAddress.Bytes())); err != nil {
		t.Error(err)
	}
	cadenceCapabilityPublicationID, _ := cadence.NewString(capabilityPublicationID)
	if err := tx.AddArgument(cadenceCapabilityPublicationID); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}

func addNftWithdrawCapAsManager(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	nftWithdrawCapStoragePathID string,
	nftCollectionStoragePathID string,
	withdrawCapTag string,
	shouldRevert bool,
) []uint64 {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, AddNftWithdrawCapAsManagerPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Add arguments to transaction
	cadenceNftWithdrawCapStoragePathID, _ := cadence.NewString(nftWithdrawCapStoragePathID)
	if err := tx.AddArgument(cadenceNftWithdrawCapStoragePathID); err != nil {
		t.Error(err)
	}
	cadenceNftCollectionStoragePathID, _ := cadence.NewString(nftCollectionStoragePathID)
	if err := tx.AddArgument(cadenceNftCollectionStoragePathID); err != nil {
		t.Error(err)
	}
	cadenceWithdrawCapTag, _ := cadence.NewString(withdrawCapTag)
	if err := tx.AddArgument(cadenceWithdrawCapTag); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)

	// Parse collection UUID from transaction result and return
	var collectionUUIDs []uint64
	if !shouldRevert {
		for _, log := range txResult.Logs {
			v, err := strconv.ParseUint(log, 10, 64)
			require.NoError(t, err)
			collectionUUIDs = append(collectionUUIDs, v)
		}
	}
	return collectionUUIDs
}

func addNftWithdrawCapAsSupplier(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	nftWithdrawCapStoragePathID string,
	nftCollectionStoragePathID string,
	withdrawCapTag string,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) []uint64 {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, AddNftWithdrawCapAsSupplierPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Add arguments to transaction
	cadenceNftWithdrawCapStoragePathID, _ := cadence.NewString(nftWithdrawCapStoragePathID)
	if err := tx.AddArgument(cadenceNftWithdrawCapStoragePathID); err != nil {
		t.Error(err)
	}
	cadenceNftCollectionStoragePathID, _ := cadence.NewString(nftCollectionStoragePathID)
	if err := tx.AddArgument(cadenceNftCollectionStoragePathID); err != nil {
		t.Error(err)
	}
	cadenceWithdrawCapTag, _ := cadence.NewString(withdrawCapTag)
	if err := tx.AddArgument(cadenceWithdrawCapTag); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)

	// Parse collection UUID from transaction result and return
	var collectionUUIDs []uint64
	if !shouldRevert {
		for _, log := range txResult.Logs {
			v, err := strconv.ParseUint(log, 10, 64)
			require.NoError(t, err)
			collectionUUIDs = append(collectionUUIDs, v)
		}
	}
	return collectionUUIDs
}

func claimAggregatedNftWithdrawCap(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	capabilityPublicationID string,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, ClaimAggregatedNftWithdrawCapPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.BytesToAddress(contracts.NFTProviderAggregatorAddress.Bytes())); err != nil {
		t.Error(err)
	}
	cadenceNftWithdrawCapStoragePathID, _ := cadence.NewString(capabilityPublicationID)
	if err := tx.AddArgument(cadenceNftWithdrawCapStoragePathID); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}

func destroyAggregator(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, DestroyAggregatorPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
}

func destroySupplier(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, DestroySupplierPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}

func publishAdditionalSupplierFactoryCap(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	supplierAddresses []flow.Address,
	capabilityPublicationIDs []string,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, PublishAdditionalSupplierFactoryCapPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Add arguments to transaction
	var cadenceSuppliers []cadence.Value
	for _, supplierAddress := range supplierAddresses {
		cadenceSuppliers = append(cadenceSuppliers, cadence.BytesToAddress(supplierAddress.Bytes()))
	}
	if err := tx.AddArgument(cadence.NewArray(cadenceSuppliers)); err != nil {
		t.Error(err)
	}
	var cadenceCapabilityPublicationIDs []cadence.Value
	for _, capabilityPublicationID := range capabilityPublicationIDs {
		cadenceCapabilityPublicationID, _ := cadence.NewString(capabilityPublicationID)
		cadenceCapabilityPublicationIDs = append(cadenceCapabilityPublicationIDs, cadenceCapabilityPublicationID)
	}
	if err := tx.AddArgument(cadence.NewArray(cadenceCapabilityPublicationIDs)); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
}

func publishAggregatedNftWithdrawCap(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipient flow.Address,
	capabilityPublicationID string,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, PublishAggregatedNftWithdrawCapPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.BytesToAddress(recipient.Bytes())); err != nil {
		t.Error(err)
	}
	cadenceCapabilityPublicationID, _ := cadence.NewString(capabilityPublicationID)
	if err := tx.AddArgument(cadenceCapabilityPublicationID); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
}

func removeNftWithdrawCapAsManager(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	collectionUUID uint64,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, RemoveNftWithdrawCapAsManagerPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.NewUInt64(collectionUUID)); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
}

func removeNftWithdrawCapAsSupplier(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	collectionUUID uint64,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, RemoveNftWithdrawCapAsSupplierPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.NewUInt64(collectionUUID)); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}

func transferFromAggregatedNftProviderAsManager(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipientAddress flow.Address,
	withdrawID uint64,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, TransferFromAggregatedNftProviderAsManagerPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTProviderAggregatorAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes())); err != nil {
		t.Error(err)
	}
	if err := tx.AddArgument(cadence.NewUInt64(withdrawID)); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTProviderAggregatorAddress},
		[]crypto.Signer{signer, contracts.NFTProviderAggregatorSigner},
		shouldRevert,
	)
}

func transferFromAggregatedNftProviderAsThirdParty(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipientAddress flow.Address,
	withdrawID uint64,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, TransferFromAggregatedNftProviderAsThirdPartyPath)).
		SetComputeLimit(200).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Add arguments to transaction
	if err := tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes())); err != nil {
		t.Error(err)
	}
	if err := tx.AddArgument(cadence.NewUInt64(withdrawID)); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}

func revokeWithdrawCapability(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	withdrawCapTag string,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// Create transaction
	tx := flow.NewTransaction().
		SetScript(loadScript(contracts, ExampleNftRevokeWithdrawCapPath)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)

	// Add arguments to transaction
	cadenceWithdrawCapTag, _ := cadence.NewString(withdrawCapTag)
	if err := tx.AddArgument(cadenceWithdrawCapTag); err != nil {
		t.Error(err)
	}

	// Sign and submit transaction
	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	return signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}
