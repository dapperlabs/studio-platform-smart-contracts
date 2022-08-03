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

func createEdition(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	metadata map[string]string,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(createEditionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.EditionNFTAddress)
	tx.AddArgument(metadataDict(metadata))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.EditionNFTAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.EditionNFTSigner},
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
		SetScript(closeEditionTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.EditionNFTAddress)
	tx.AddArgument(cadence.NewUInt64(editionID))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.EditionNFTAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.EditionNFTSigner},
		shouldRevert,
	)
}

func mintSeasonalNFT(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	recipientAddress flow.Address,
	editionID uint64,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(mintEditionNFTTransaction(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(contracts.EditionNFTAddress)
	tx.AddArgument(cadence.BytesToAddress(recipientAddress.Bytes()))
	tx.AddArgument(cadence.NewUInt64(editionID))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, contracts.EditionNFTAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), contracts.EditionNFTSigner},
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
		SetScript(loadEditionNFTTransferNFTTransaction(contracts)).
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
