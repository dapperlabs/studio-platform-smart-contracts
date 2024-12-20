package test

import (
	"regexp"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftAddressPlaceholder          = "\"NonFungibleToken\""
	ftAddressPlaceholder           = "\"FungibleToken\""
	mvAddressPlaceholder           = "\"MetadataViews\""
	viewResolverAddressPlaceholder = "\"ViewResolver\""
	AllDayAddressPlaceholder       = "\"AllDay\""
	royaltyAddressPlaceholder      = "0xALLDAYROYALTYADDRESS"
	escrowAddressPlaceholder       = "\"Escrow\""
	nftLockerAddressPlaceholder    = "\"NFTLocker\""
	escrowAddressPlaceholderBis    = "0xf8d6e0586b0a20c7"
)

const (
	EscrowPath                 = "../../../contracts/Escrow.cdc"
	AllDayPath                 = "../../../contracts/AllDay.cdc"
	NFTLockerPath              = "../../../../locked-nft/contracts/NFTLocker.cdc"
	EscrowTransactionsRootPath = "../../../transactions"
	EscrowScriptsRootPath      = "../../../scripts"

	// Accounts
	AllDaySetupAccountPath   = EscrowTransactionsRootPath + "/user/setup_allday_account.cdc"
	AllDayAccountIsSetupPath = EscrowScriptsRootPath + "/user/account_is_setup.cdc"

	// Series
	EscrowCreateSeriesPath   = EscrowTransactionsRootPath + "/admin/series/create_series.cdc"
	EscrowReadSeriesByIDPath = EscrowScriptsRootPath + "/series/read_series_by_id.cdc"

	// Sets
	EscrowCreateSetPath   = EscrowTransactionsRootPath + "/admin/sets/create_set.cdc"
	EscrowReadSetByIDPath = EscrowScriptsRootPath + "/sets/read_set_by_id.cdc"

	// Plays
	EscrowCreatePlayPath   = EscrowTransactionsRootPath + "/admin/plays/create_play.cdc"
	EscrowReadPlayByIDPath = EscrowScriptsRootPath + "/plays/read_play_by_id.cdc"

	// Editions
	EscrowCreateEditionPath   = EscrowTransactionsRootPath + "/admin/editions/create_edition.cdc"
	EscrowReadEditionByIDPath = EscrowScriptsRootPath + "/editions/read_edition_by_id.cdc"

	// Leaderboards
	EscrowCreateLeaderboardPath = EscrowTransactionsRootPath + "/admin/leaderboards/create_leaderboard.cdc"

	// Moment NFTs
	EscrowMintMomentNFTPath           = EscrowTransactionsRootPath + "/admin/nfts/mint_moment_nft.cdc"
	EscrowReadMomentNFTSupplyPath     = EscrowScriptsRootPath + "/nfts/read_moment_nft_supply.cdc"
	EscrowReadMomentNFTPropertiesPath = EscrowScriptsRootPath + "/nfts/read_moment_nft_properties.cdc"
	EscrowReadCollectionLengthPath    = EscrowScriptsRootPath + "/nfts/read_collection_nft_length.cdc"

	// Escrow
	EscrowMomentNFTPath              = EscrowTransactionsRootPath + "/user/add_entry.cdc"
	EscrowWithdrawMomentNFTPath      = EscrowTransactionsRootPath + "/admin/leaderboards/withdraw_entry.cdc"
	EscrowAdminTransferMomentNFTPath = EscrowTransactionsRootPath + "/admin/leaderboards/admin_transfer.cdc"
	EscrowBurnNFTPath                = EscrowTransactionsRootPath + "/admin/leaderboards/burn_nft.cdc"
	EscrowReadLeaderboardInfoPath    = EscrowScriptsRootPath + "/leaderboards/read_leaderboard_info.cdc"
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	ftRe := regexp.MustCompile(ftAddressPlaceholder)
	code = ftRe.ReplaceAll(code, []byte("0x"+ftAddress.String()))

	AllDayRe := regexp.MustCompile(AllDayAddressPlaceholder)
	code = AllDayRe.ReplaceAll(code, []byte("0x"+contracts.AllDayAddress.String()))

	mvRe := regexp.MustCompile(mvAddressPlaceholder)
	code = mvRe.ReplaceAll(code, []byte("0x"+contracts.MetadataViewsAddress.String()))

	royaltyRe := regexp.MustCompile(royaltyAddressPlaceholder)
	code = royaltyRe.ReplaceAll(code, []byte("0x"+contracts.RoyaltyAddress.String()))

	escrowRe := regexp.MustCompile(escrowAddressPlaceholder)
	code = escrowRe.ReplaceAll(code, []byte("0x"+contracts.EscrowAddress.String()))

	escrowReBis := regexp.MustCompile(escrowAddressPlaceholderBis)
	code = escrowReBis.ReplaceAll(code, []byte("0x"+contracts.EscrowAddress.String()))

	return code
}

func LoadAllDay(nftAddress, metaAddress, viewResolverAddress, royaltyAddress flow.Address) []byte {
	code := readFile(AllDayPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	ftRe := regexp.MustCompile(ftAddressPlaceholder)
	code = ftRe.ReplaceAll(code, []byte("0x"+ftAddress.String()))

	mvRe := regexp.MustCompile(mvAddressPlaceholder)
	code = mvRe.ReplaceAll(code, []byte("0x"+metaAddress.String()))

	vRe := regexp.MustCompile(viewResolverAddressPlaceholder)
	code = vRe.ReplaceAll(code, []byte("0x"+viewResolverAddress.String()))

	royaltyRe := regexp.MustCompile(royaltyAddressPlaceholder)
	code = royaltyRe.ReplaceAll(code, []byte("0x"+royaltyAddress.String()))

	return code
}

func LoadEscrow(nftAddress, nftLockerAddress flow.Address) []byte {
	code := readFile(EscrowPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	nftLockerRe := regexp.MustCompile(nftLockerAddressPlaceholder)
	code = nftLockerRe.ReplaceAll(code, []byte("0x"+nftLockerAddress.String()))

	return code
}

func LoadNFTLockerContract(nftAddress flow.Address) []byte {
	code := readFile(NFTLockerPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	return code
}

func loadAllDaySetupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDaySetupAccountPath),
		contracts,
	)
}

func loadAllDayAccountIsSetupScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(AllDayAccountIsSetupPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Series
// ------------------------------------------------------------
func loadEscrowCreateSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowCreateSeriesPath),
		contracts,
	)
}

func loadEscrowReadSeriesByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadSeriesByIDPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Sets
// ------------------------------------------------------------
func loadEscrowCreateSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowCreateSetPath),
		contracts,
	)
}

func loadEscrowReadSetByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadSetByIDPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Plays
// ------------------------------------------------------------
func loadEscrowCreatePlayTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowCreatePlayPath),
		contracts,
	)
}

func loadEscrowReadPlayByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadPlayByIDPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Editions
// ------------------------------------------------------------
func loadEscrowCreateEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowCreateEditionPath),
		contracts,
	)
}

func loadEscrowReadEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadEditionByIDPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Moment NFTs
// ------------------------------------------------------------
func loadEscrowMintMomentNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowMintMomentNFTPath),
		contracts,
	)
}

func loadEscrowReadCollectionLengthScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadCollectionLengthPath),
		contracts,
	)
}

func loadEscrowReadMomentNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadMomentNFTSupplyPath),
		contracts,
	)
}

func loadEscrowReadMomentNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadMomentNFTPropertiesPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Escrow
// ------------------------------------------------------------
func loadEscrowMomentNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowMomentNFTPath),
		contracts,
	)
}

func loadEscrowLeaderboardInfoScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowReadLeaderboardInfoPath),
		contracts,
	)
}

func loadEscrowWithdrawMomentNFT(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowWithdrawMomentNFTPath),
		contracts,
	)
}

func loadEscrowAdminTransferMomentNFT(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowAdminTransferMomentNFTPath),
		contracts,
	)
}

func loadEscrowBurnNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowBurnNFTPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Leaderboards
// ------------------------------------------------------------
func loadEscrowCreateLeaderboardTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EscrowCreateLeaderboardPath),
		contracts,
	)
}
