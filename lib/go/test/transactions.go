package test

import (
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	fttemplates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
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

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address},
		[]crypto.Signer{b.ServiceKey().Signer()},
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
		SetScript(loadSportCreateSeriesTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadenceString)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportCloseSeriesTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadence.NewUInt64(id))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportCreateSetTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadenceString)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportCreatePlayTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadenceString)
	tx.AddArgument(metadataDict(metadata))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportCreateEditionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadence.NewUInt64(seriesID))
	tx.AddArgument(cadence.NewUInt64(setID))
	tx.AddArgument(cadence.NewUInt64(playID))
	tx.AddArgument(cadenceString)
	if maxMintSize != nil {
		tx.AddArgument(cadence.NewUInt64(*maxMintSize))
	} else {
		tx.AddArgument(cadence.Optional{})
	}

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportCloseEditionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadence.NewUInt64(editionID))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportMintMomentNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.SportAddress)
	tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes()))
	tx.AddArgument(cadence.NewUInt64(editionID))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.SportAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.SportSigner},
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
		SetScript(loadSportTransferNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)
	tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes()))
	tx.AddArgument(cadence.NewUInt64(nftID))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), senderSigner},
		shouldRevert,
	)
}
