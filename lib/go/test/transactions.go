package test

import (
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	fttemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/require"
)

//------------------------------------------------------------
// Setup
//------------------------------------------------------------
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

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address},
		[]crypto.Signer{signer},
		false,
	)
}

//------------------------------------------------------------
// Series
//------------------------------------------------------------
func createSeries(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	name string,
	shouldRevert bool,
) {
	cadenceString, _ := cadence.NewString(name)
	tx := flow.NewTransaction().
		SetScript(loadDapperSportCreateSeriesTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadenceString)

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

func closeSeries(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	id uint64,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(loadDapperSportCloseSeriesTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadence.NewUInt64(id))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

//------------------------------------------------------------
// Sets
//------------------------------------------------------------
func createSet(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	name string,
	shouldRevert bool,
) {
	cadenceString, _ := cadence.NewString(name)
	tx := flow.NewTransaction().
		SetScript(loadDapperSportCreateSetTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadenceString)

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

//------------------------------------------------------------
// Plays
//------------------------------------------------------------
func createPlay(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	classification string,
	metadata map[string]string,
	shouldRevert bool,
) {
	cadenceString, _ := cadence.NewString(classification)
	tx := flow.NewTransaction().
		SetScript(loadDapperSportCreatePlayTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadenceString)
	tx.AddArgument(metadataDict(metadata))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

//------------------------------------------------------------
// Editions
//------------------------------------------------------------
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
	cadenceString, _ := cadence.NewString(tier)
	tx := flow.NewTransaction().
		SetScript(loadDapperSportCreateEditionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadence.NewUInt64(seriesID))
	tx.AddArgument(cadence.NewUInt64(setID))
	tx.AddArgument(cadence.NewUInt64(playID))
	tx.AddArgument(cadenceString)
	if maxMintSize != nil {
		tx.AddArgument(cadence.NewUInt64(*maxMintSize))
	} else {
		tx.AddArgument(cadence.Optional{})
	}

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

func closeEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	editionID uint64,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(loadDapperSportCloseEditionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadence.NewUInt64(editionID))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

//------------------------------------------------------------
// MomentNFTs
//------------------------------------------------------------
func mintMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipientAddress flow.Address,
	editionID uint64,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(loadDapperSportMintMomentNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.DapperSportAddress)
	tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes()))
	tx.AddArgument(cadence.NewUInt64(editionID))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.DapperSportAddress},
		[]crypto.Signer{signer, contracts.DapperSportSigner},
		shouldRevert,
	)
}

func transferMomentNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	senderAddress flow.Address,
	senderSigner crypto.Signer,
	nftID uint64,
	recipientAddress flow.Address,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(loadDapperSportTransferNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)
	tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes()))
	tx.AddArgument(cadence.NewUInt64(nftID))

	signer, err := b.ServiceKey().Signer()
	require.NoError(t, err)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{signer, senderSigner},
		shouldRevert,
	)
}
