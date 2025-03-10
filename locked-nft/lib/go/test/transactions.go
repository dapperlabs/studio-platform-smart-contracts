package test

import (
	"strings"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/require"
)

// ------------------------------------------------------------
// Setup
// ------------------------------------------------------------

func mintExampleNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	recipientAddress string,
) uint64 {
	tx := flow.NewTransaction().
		SetScript(mintExampleNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTLockerAddress)
	tx.AddArgument(cadence.Address(flow.HexToAddress(recipientAddress)))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTLockerAddress},
		[]crypto.Signer{signer, contracts.NFTLockerSigner},
		shouldRevert,
	)

	nftId := uint64(0)
	for _, event := range txResult.Events {
		if strings.Contains(event.Type, "NonFungibleToken.Deposited") {
			if v := cadence.FieldsMappedByName(event.Value)["id"]; v != nil {
				nftId = GetFieldValue(v).(uint64)
			}
		}
	}

	return nftId
}

func lockNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftId uint64,
	duration uint64,
) (uint64, uint64) {
	tx := flow.NewTransaction().
		SetScript(lockNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)
	tx.AddArgument(cadence.UInt64(nftId))
	tx.AddArgument(cadence.UInt64(duration))

	signer, _ := b.ServiceKey().Signer()
	txResult := signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{signer, userSigner},
		shouldRevert,
	)

	var lockedAt, lockedUntil uint64

	for _, event := range txResult.Events {
		if strings.Contains(event.Type, "NFTLocker.NFTLocked") {
			if v := cadence.FieldsMappedByName(event.Value)["lockedAt"]; v != nil {
				lockedAt = GetFieldValue(v).(uint64)
			}
			if v := cadence.FieldsMappedByName(event.Value)["lockedUntil"]; v != nil {
				lockedUntil = GetFieldValue(v).(uint64)
			}
			break
		}
	}

	return lockedAt, lockedUntil
}

func unlockNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftId uint64,
) {
	tx := flow.NewTransaction().
		SetScript(unlockNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)
	tx.AddArgument(cadence.UInt64(nftId))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{signer, userSigner},
		shouldRevert,
	)
}

func adminAddReceiver(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(adminAddReceiverTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTLockerAddress)

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTLockerAddress},
		[]crypto.Signer{signer, contracts.NFTLockerSigner},
		shouldRevert,
	)
}

func adminRemoveReceiver(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(adminRemoveReceiverTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTLockerAddress)

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTLockerAddress},
		[]crypto.Signer{signer, contracts.NFTLockerSigner},
		shouldRevert,
	)
}

func unlockNFTWithAuthorizedDeposit(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	userAddress flow.Address,
	userSigner crypto.Signer,
	leaderboardName string,
	nftId uint64,
) {
	tx := flow.NewTransaction().
		SetScript(unlockNFTWithAuthorizedDepositTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	leaderboardNameCadence, _ := cadence.NewString(leaderboardName)
	tx.AddArgument(leaderboardNameCadence)
	tx.AddArgument(cadence.UInt64(nftId))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{signer, userSigner},
		shouldRevert,
	)
}

func adminUnlockNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	shouldRevert bool,
	nftId uint64,
) {
	tx := flow.NewTransaction().
		SetScript(adminUnlockNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTLockerAddress)
	tx.AddArgument(cadence.UInt64(nftId))

	signer, _ := b.ServiceKey().Signer()
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTLockerAddress},
		[]crypto.Signer{signer, contracts.NFTLockerSigner},
		shouldRevert,
	)
}

func createLeaderboard(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	leaderboardName string,
) {
	tx := flow.NewTransaction().
		SetScript(createLeaderboardTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.NFTLockerAddress)
	tx.AddArgument(cadence.String(leaderboardName))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.NFTLockerAddress},
		[]crypto.Signer{signer, contracts.NFTLockerSigner},
		false,
	)
}
