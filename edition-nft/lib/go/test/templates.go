package test

import (
	"regexp"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftAddressPlaceholder        = "\"[^\"]*NonFungibleToken.cdc\""
	EditionNFTAddressPlaceholder = "\"[^\"]*EditionNFT.cdc\""
)

const (
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
	AllDayTransferNFTPath   = TransactionsRootPath + "/user/transfer_moment_nft.cdc"
	ReadNftSupplyScriptPath = ScriptsRootPath + "/nfts/read_nft_supply.cdc"
	ReadNftPropertiesTxPath = ScriptsRootPath + "/nfts/read_nft_properties.cdc"
)

//------------------------------------------------------------
// Accounts
//------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	EditionNFTRe := regexp.MustCompile(EditionNFTAddressPlaceholder)
	code = EditionNFTRe.ReplaceAll(code, []byte("0x"+contracts.AllDayAddress.String()))

	return code
}

func allDaySeasonalContract(nftAddress flow.Address) []byte {
	code := readFile(EditionNFTPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

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

//------------------------------------------------------------
// Moment NFTs
//------------------------------------------------------------
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

func loadAllDayTransferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDayTransferNFTPath),
		contracts,
	)
}
