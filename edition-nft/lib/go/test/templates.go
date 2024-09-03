package test

import (
	"strings"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftAddressPlaceholder           = "\"NonFungibleToken\""
	metadataViewsAddressPlaceholder = "\"MetadataViews\""
	viewResolverAddressPlaceholder  = "\"ViewResolver\""
	EditionNFTAddressPlaceholder    = "\"EditionNFT\""
)

const (
	BurnerPath           = "../../../contracts/imports/Burner.cdc"
	EditionNFTPath       = "../../../contracts/EditionNFT.cdc"
	TransactionsRootPath = "../../../transactions"
	ScriptsRootPath      = "../../../scripts"

	// Accounts
	SetupAccountTxPath       = TransactionsRootPath + "/user/setup_user_account.cdc"
	IsAccountSetupScriptPath = ScriptsRootPath + "/user/is_account_setup.cdc"

	// Editions
	CreateEditionTxPath       = TransactionsRootPath + "/admin/editions/create_edition.cdc"
	CloseEditionTxPath        = TransactionsRootPath + "/admin/editions/close_edition.cdc"
	ReadEditionByIDScriptPath = ScriptsRootPath + "/editions/read_edition_by_id.cdc"

	// NFTs
	MintNFTTxPath           = TransactionsRootPath + "/admin/nfts/mint_nft.cdc"
	ReadNftSupplyScriptPath = ScriptsRootPath + "/nfts/read_nft_supply.cdc"
	ReadNftPropertiesTxPath = ScriptsRootPath + "/nfts/read_nft_properties.cdc"
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	code = []byte(strings.ReplaceAll(string(code), nftAddressPlaceholder, "0x"+contracts.NFTAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewsAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), viewResolverAddressPlaceholder, "0x"+contracts.ViewResolverAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), EditionNFTAddressPlaceholder, "0x"+contracts.EditionNFTAddress.String()))

	return code
}

func allDaySeasonalContract(nftAddress flow.Address, metadataViewsAddr flow.Address, viewResolverAddress flow.Address) []byte {
	code := readFile(EditionNFTPath)

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddr.String()))
	code = []byte(strings.ReplaceAll(string(code), nftAddressPlaceholder, "0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), viewResolverAddressPlaceholder, "0x"+viewResolverAddress.String()))

	return code
}

func setupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SetupAccountTxPath),
		contracts,
	)
}

func isAccountSetupScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(IsAccountSetupScriptPath),
		contracts,
	)
}

func createEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CreateEditionTxPath),
		contracts,
	)
}

func readEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(ReadEditionByIDScriptPath),
		contracts,
	)
}

func closeEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CloseEditionTxPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Moment NFTs
// ------------------------------------------------------------
func mintEditionNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(MintNFTTxPath),
		contracts,
	)
}

func getEditionNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(ReadNftSupplyScriptPath),
		contracts,
	)
}

func getEditionNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(ReadNftPropertiesTxPath),
		contracts,
	)
}
