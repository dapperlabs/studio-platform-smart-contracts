package test

import (
	"regexp"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftAddressPlaceholder            = "\"[^\"]*NonFungibleToken.cdc\""
	AllDayAddressPlaceholder         = "\"[^\"]*AllDay.cdc\""
	AllDaySeasonalAddressPlaceholder = "\"[^\"]*AllDaySeasonal.cdc\""
)

const (
	AllDayPath                 = "../../../contracts/AllDay.cdc"
	AllDaySeasonalPath         = "../../../contracts/AllDaySeasonal.cdc"
	AllDayTransactionsRootPath = "../../../transactions"
	AllDayScriptsRootPath      = "../../../scripts"

	// Accounts
	AllDaySetupAccountPath         = AllDayTransactionsRootPath + "/user/setup_allday_account.cdc"
	AllDaySeasonalSetupAccountPath = AllDayTransactionsRootPath + "/user/setup_allday_seasonal_account.cdc"
	IsAccountSetupScriptPath       = AllDayScriptsRootPath + "/user/is_account_setup.cdc"

	// Editions
	AllDayCreateEditionPath         = AllDayTransactionsRootPath + "/admin/editions/create_edition.cdc"
	AllDaySeasonalCreateEditionPath = AllDayTransactionsRootPath + "/admin/editions/create_seasonal_edition.cdc"
	AllDayCloseEditionPath          = AllDayTransactionsRootPath + "/admin/editions/close_edition.cdc"
	AllDaySeasonalCloseEditionPath  = AllDayTransactionsRootPath + "/admin/editions/close_seasonal_edition.cdc"

	AllDayReadEditionByIDPath         = AllDayScriptsRootPath + "/editions/read_edition_by_id.cdc"
	AllDaySeasonalReadEditionByIDPath = AllDayScriptsRootPath + "/editions/read_seasonal_edition_by_id.cdc"
	AllDayReadAllEditionsPath         = AllDayScriptsRootPath + "/edition/read_all_editions.cdc"

	// Moment NFTs
	AllDayMintMomentNFTPath   = AllDayTransactionsRootPath + "/admin/nfts/mint_moment_nft.cdc"
	AllDayMintSeasonalNFTPath = AllDayTransactionsRootPath + "/admin/nfts/mint_seasonal_nft.cdc"

	AllDayMintMomentNFTMultiPath        = AllDayTransactionsRootPath + "/admin/nfts/mint_moment_nft_multi.cdc"
	AllDayTransferNFTPath               = AllDayTransactionsRootPath + "/user/transfer_moment_nft.cdc"
	AllDayReadMomentNFTSupplyPath       = AllDayScriptsRootPath + "/nfts/read_moment_nft_supply.cdc"
	AllDayReadSeasonalNFTSupplyPath     = AllDayScriptsRootPath + "/nfts/read_seasonal_nft_supply.cdc"
	AllDayReadMomentNFTPropertiesPath   = AllDayScriptsRootPath + "/nfts/read_moment_nft_properties.cdc"
	AllDayReadSeasonalNFTPropertiesPath = AllDayScriptsRootPath + "/nfts/read_seasonal_nft_properties.cdc"

	AllDayReadCollectionNFTLengthPath = AllDayScriptsRootPath + "/nfts/read_collection_nft_length.cdc"
	AllDayReadCollectionNFTIDsPath    = AllDayScriptsRootPath + "/nfts/read_collection_nft_ids.cdc"
)

//------------------------------------------------------------
// Accounts
//------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	AllDayRe := regexp.MustCompile(AllDayAddressPlaceholder)
	code = AllDayRe.ReplaceAll(code, []byte("0x"+contracts.AllDayAddress.String()))

	AllDaySeasonalRe := regexp.MustCompile(AllDaySeasonalAddressPlaceholder)
	code = AllDaySeasonalRe.ReplaceAll(code, []byte("0x"+contracts.AllDayAddress.String()))

	return code
}

func allDaySeasonalContract(nftAddress flow.Address) []byte {
	code := readFile(AllDaySeasonalPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	return code
}

func setupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDaySeasonalSetupAccountPath),
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
		readFile(AllDaySeasonalCreateEditionPath),
		contracts,
	)
}

func readEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDaySeasonalReadEditionByIDPath),
		contracts,
	)
}

func closeEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDaySeasonalCloseEditionPath),
		contracts,
	)
}

//func loadAllDayReadAllEditionsScript(contracts Contracts) []byte {
//	return replaceAddresses(
//		readFile(AllDayReadAllEditionsPath),
//		contracts,
//	)
//}

//------------------------------------------------------------
// Moment NFTs
//------------------------------------------------------------
func mintEditionNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDayMintSeasonalNFTPath),
		contracts,
	)
}

func getEditionNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDayReadSeasonalNFTSupplyPath),
		contracts,
	)
}

func getEditionNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDayReadSeasonalNFTPropertiesPath),
		contracts,
	)
}

//func readCollectionNFTLengthScript(contracts Contracts) []byte {
//	return replaceAddresses(
//		readFile(AllDayReadCollectionNFTLengthPath),
//		contracts,
//	)
//}
//
//func loadAllDayReadCollectionNFTIDsScript(contracts Contracts) []byte {
//	return replaceAddresses(
//		readFile(AllDayReadCollectionNFTIDsPath),
//		contracts,
//	)
//}

func loadAllDayTransferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDayTransferNFTPath),
		contracts,
	)
}
