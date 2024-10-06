package test

import (
	"math/big"
	"testing"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

// ------------------------------------------------------------
// Series
// ------------------------------------------------------------
func createSeries(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	name string,
	shouldRevert bool,
) {
	nameString, _ := cadence.NewString(name)
	tx := flow.NewTransaction().
		SetScript(loadEscrowCreateSeriesTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(nameString)

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		shouldRevert,
	)
}

// ------------------------------------------------------------
// Sets
// ------------------------------------------------------------
func createSet(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	name string,
	shouldRevert bool,
) {
	nameString, _ := cadence.NewString(name)
	tx := flow.NewTransaction().
		SetScript(loadEscrowCreateSetTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(nameString)

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		shouldRevert,
	)
}

// ------------------------------------------------------------
// Plays
// ------------------------------------------------------------
func createPlay(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	classification string,
	metadata map[string]string,
	shouldRevert bool,
) {
	classificationString, _ := cadence.NewString(classification)
	tx := flow.NewTransaction().
		SetScript(loadEscrowCreatePlayTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(classificationString)
	tx.AddArgument(metadataDict(metadata))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		shouldRevert,
	)
}

// ------------------------------------------------------------
// Editions
// ------------------------------------------------------------
func createEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	seriesID uint64,
	setID uint64,
	playID uint64,
	maxMintSize *uint64,
	tier string,
	shouldRevert bool,
) {
	tierString, _ := cadence.NewString(tier)
	tx := flow.NewTransaction().
		SetScript(loadEscrowCreateEditionTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(cadence.NewUInt64(seriesID))
	tx.AddArgument(cadence.NewUInt64(setID))
	tx.AddArgument(cadence.NewUInt64(playID))
	tx.AddArgument(tierString)
	if maxMintSize != nil {
		tx.AddArgument(cadence.NewUInt64(*maxMintSize))
	} else {
		tx.AddArgument(cadence.Optional{})
	}

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		shouldRevert,
	)
}

// ------------------------------------------------------------
// Leaderboards
// ------------------------------------------------------------
func createLeaderboard(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	leaderboardName string,
) {
	tx := flow.NewTransaction().
		SetScript(loadEscrowCreateLeaderboardTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(cadence.String(leaderboardName))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		false,
	)
}

func getLeaderboardData(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	leaderboardName string,
) (LeaderboardInfo, error) {
	script := loadEscrowLeaderboardInfoScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{
		jsoncdc.MustEncode(cadence.String(leaderboardName))})

	return parseLeaderboardInfo(result)
}

// ------------------------------------------------------------
// MomentNFTs
// ------------------------------------------------------------
func escrowMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	userSigner crypto.Signer,
	ownerAddress flow.Address,
	momentNftFlowID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(loadEscrowMomentNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(ownerAddress)
	tx.AddArgument(cadence.String("leaderboardBurn-1"))
	tx.AddArgument(cadence.NewUInt64(momentNftFlowID))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, ownerAddress},
		[]crypto.Signer{signer, userSigner},
		false,
	)
}

func getMomentNFTLengthInAccount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	address flow.Address,
) *big.Int {
	script := loadEscrowReadCollectionLengthScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.BytesToAddress(address.Bytes()))})

	return GetFieldValue(result).(*big.Int)
}

func mintMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipientAddress flow.Address,
	editionID uint64,
	serialNumber *uint64,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(loadEscrowMintMomentNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes()))
	tx.AddArgument(cadence.NewUInt64(editionID))
	sNumber := cadence.NewOptional(nil)
	if serialNumber != nil {
		sNumber = cadence.NewOptional(cadence.NewUInt64(*serialNumber))
	}
	tx.AddArgument(sNumber)

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		shouldRevert,
	)
}

// ------------------------------------------------------------
// Escrow
// ------------------------------------------------------------
func withdrawMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	leaderboardName string,
	userAddress flow.Address,
	momentNftFlowID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(loadEscrowWithdrawMomentNFT(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(cadence.String(leaderboardName))
	tx.AddArgument(cadence.NewUInt64(momentNftFlowID))
	tx.AddArgument(cadence.BytesToAddress(userAddress.Bytes()))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		false,
	)
}

func adminTransferMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	leaderboardName string,
	userAddress flow.Address,
	momentNftFlowID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(loadEscrowAdminTransferMomentNFT(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(cadence.String(leaderboardName))
	tx.AddArgument(cadence.NewUInt64(momentNftFlowID))
	tx.AddArgument(cadence.BytesToAddress(userAddress.Bytes()))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		false,
	)
}

func burnMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	leaderboardName string,
	momentNftFlowID uint64,
) {
	tx := flow.NewTransaction().
		SetScript(loadEscrowBurnNFTTransaction(contracts)).
		SetComputeLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.AllDayAddress)
	tx.AddArgument(cadence.String(leaderboardName))
	tx.AddArgument(cadence.NewUInt64(momentNftFlowID))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)
	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.AllDayAddress},
		[]crypto.Signer{signer, contracts.AllDaySigner},
		false,
	)
}
