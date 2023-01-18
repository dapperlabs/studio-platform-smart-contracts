package test

import (
	"regexp"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftAddressPlaceholder           = "\"[^\"]*NonFungibleToken.cdc\""
	DSSCollectionAddressPlaceholder = "\"[^\"]*DSSCollection.cdc\""
)

const (
	DSSCollectionPath    = "../../../contracts/DSSCollection.cdc"
	TransactionsRootPath = "../../../transactions"
	ScriptsRootPath      = "../../../scripts"

	// Accounts
	SetupAccountTxPath       = TransactionsRootPath + "/user/setup_user_account.cdc"
	IsAccountSetupScriptPath = ScriptsRootPath + "/user/is_account_setup.cdc"

	// Collection Groups
	CreateCollectionGroupTxPath          = TransactionsRootPath + "/admin/create_collection_group.cdc"
	CreateTimeBoundCollectionGroupTxPath = TransactionsRootPath + "/admin/create_collection_group_time_bound.cdc"
	CloseCollectionGroupTxPath           = TransactionsRootPath + "/admin/close_collection_group.cdc"
	AddNFTToCollectionGroupTxPath        = TransactionsRootPath + "/admin/add_nft_to_collection_group.cdc"
	GetCollectionGroupByIDScriptPath     = ScriptsRootPath + "/get_collection_group.cdc"

	// NFTs
	MintNFTTxPath           = TransactionsRootPath + "/admin/mint_nft.cdc"
	ReadNftSupplyScriptPath = ScriptsRootPath + "/total_supply.cdc"
	ReadNftPropertiesTxPath = ScriptsRootPath + "/get_nft.cdc"
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	DSSCollectionRe := regexp.MustCompile(DSSCollectionAddressPlaceholder)
	code = DSSCollectionRe.ReplaceAll(code, []byte("0x"+contracts.DSSCollectionAddress.String()))

	return code
}

func dssCollectionContract(nftAddress flow.Address) []byte {
	code := readFile(DSSCollectionPath)

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

func createCollectionGroupTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CreateCollectionGroupTxPath),
		contracts,
	)
}

func createTimeBoundCollectionGroupTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CreateTimeBoundCollectionGroupTxPath),
		contracts,
	)
}

func readCollectionGroupByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GetCollectionGroupByIDScriptPath),
		contracts,
	)
}

func closeCollectionGroupTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CloseCollectionGroupTxPath),
		contracts,
	)
}

func addNFTToCollectionGroupTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AddNFTToCollectionGroupTxPath),
		contracts,
	)
}

// ------------------------------------------------------------
// DSSCollection NFTs
// ------------------------------------------------------------
func mintDSSCollectionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(MintNFTTxPath),
		contracts,
	)
}

func readDSSCollectionSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(ReadNftSupplyScriptPath),
		contracts,
	)
}

func readNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(ReadNftPropertiesTxPath),
		contracts,
	)
}
