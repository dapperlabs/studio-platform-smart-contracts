package test

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"strings"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftAddressPlaceholder           = "\"[^\"]*NonFungibleToken.cdc\""
	SportAddressPlaceholder         = "\"[^\"]*Sport.cdc\""
	metadataViewsAddressPlaceholder = "0xMETADATAVIEWSADDRESS"

	SportPath                 = "../../../contracts/Sport.cdc"
	SportTransactionsRootPath = "../../../transactions"
	SportScriptsRootPath      = "../../../scripts"

	// Accounts
	SportSetupAccountPath   = SportTransactionsRootPath + "/user/setup_Sport_account.cdc"
	SportAccountIsSetupPath = SportScriptsRootPath + "/user/account_is_setup.cdc"

	// Series
	SportCreateSeriesPath       = SportTransactionsRootPath + "/admin/series/create_series.cdc"
	SportCloseSeriesPath        = SportTransactionsRootPath + "/admin/series/close_series.cdc"
	SportReadAllSeriesPath      = SportScriptsRootPath + "/series/read_all_series.cdc"
	SportReadSeriesByIDPath     = SportScriptsRootPath + "/series/read_series_by_id.cdc"
	SportReadSeriesByNamePath   = SportScriptsRootPath + "/series/read_series_by_name.cdc"
	SportReadAllSeriesNamesPath = SportScriptsRootPath + "/series/read_all_series_names.cdc"

	// Sets
	SportCreateSetPath       = SportTransactionsRootPath + "/admin/sets/create_set.cdc"
	SportReadAllSetsPath     = SportScriptsRootPath + "/sets/read_all_sets.cdc"
	SportReadSetByIDPath     = SportScriptsRootPath + "/sets/read_set_by_id.cdc"
	SportReadSetsByNamePath  = SportScriptsRootPath + "/sets/read_sets_by_name.cdc"
	SportReadAllSetNamesPath = SportScriptsRootPath + "/sets/read_all_set_names.cdc"

	// Plays
	SportCreatePlayPath   = SportTransactionsRootPath + "/admin/plays/create_play.cdc"
	SportReadPlayByIDPath = SportScriptsRootPath + "/plays/read_play_by_id.cdc"
	SportReadAllPlaysPath = SportScriptsRootPath + "/plays/read_all_plays.cdc"

	// Editions
	SportCreateEditionPath   = SportTransactionsRootPath + "/admin/editions/create_edition.cdc"
	SportCloseEditionPath    = SportTransactionsRootPath + "/admin/editions/close_edition.cdc"
	SportReadEditionByIDPath = SportScriptsRootPath + "/editions/read_edition_by_id.cdc"
	SportReadAllEditionsPath = SportScriptsRootPath + "/edition/read_all_editions.cdc"

	// Moment NFTs
	SportMintMomentNFTPath           = SportTransactionsRootPath + "/admin/nfts/mint_moment_nft.cdc"
	SportMintMomentNFTMultiPath      = SportTransactionsRootPath + "/admin/nfts/mint_moment_nft_multi.cdc"
	SportTransferNFTPath             = SportTransactionsRootPath + "/user/transfer_moment_nft.cdc"
	SportReadMomentNFTSupplyPath     = SportScriptsRootPath + "/nfts/read_moment_nft_supply.cdc"
	SportReadMomentNFTPropertiesPath = SportScriptsRootPath + "/nfts/read_moment_nft_properties.cdc"
	SportReadCollectionNFTLengthPath = SportScriptsRootPath + "/nfts/read_collection_nft_length.cdc"
	SportReadCollectionNFTIDsPath    = SportScriptsRootPath + "/nfts/read_collection_nft_ids.cdc"
	SportDisplayMetadataViewPath     = SportScriptsRootPath + "/nfts/metadata_display_view.cdc"

	// MetadataViews
	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"
	MetadataFTReplaceAddress      = `"./utility/FungibleToken.cdc"`
	MetadataNFTReplaceAddress     = `"./NonFungibleToken.cdc"`
)

//------------------------------------------------------------
// Accounts
//------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	SportRe := regexp.MustCompile(SportAddressPlaceholder)
	code = SportRe.ReplaceAll(code, []byte("0x"+contracts.SportAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewAddress.String()))

	return code
}

func DownloadFile(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func LoadMetadataViews(ftAddress flow.Address, nftAddress flow.Address) []byte {
	code, _ := DownloadFile(MetadataViewsContractsBaseURL + MetadataViewsInterfaceFile)
	code = []byte(strings.Replace(strings.Replace(string(code), MetadataFTReplaceAddress, "0x"+ftAddress.String(), 1), MetadataNFTReplaceAddress, "0x"+nftAddress.String(), 1))

	return code
}

func LoadSport(nftAddress flow.Address, metadataViewsAddr flow.Address) []byte {
	code := readFile(SportPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddr.String()))

	return code
}

func loadSportSetupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportSetupAccountPath),
		contracts,
	)
}

func loadSportAccountIsSetupScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportAccountIsSetupPath),
		contracts,
	)
}

//------------------------------------------------------------
// Series
//------------------------------------------------------------
func loadSportCreateSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportCreateSeriesPath),
		contracts,
	)
}

func loadSportReadSeriesByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadSeriesByIDPath),
		contracts,
	)
}

func loadSportReadSeriesByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadSeriesByNamePath),
		contracts,
	)
}

func loadSportReadAllSeriesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadAllSeriesPath),
		contracts,
	)
}

func loadSportReadAllSeriesNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadAllSeriesNamesPath),
		contracts,
	)
}

func loadSportCloseSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportCloseSeriesPath),
		contracts,
	)
}

//------------------------------------------------------------
// Sets
//------------------------------------------------------------
func loadSportCreateSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportCreateSetPath),
		contracts,
	)
}

func loadSportReadSetByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadSetByIDPath),
		contracts,
	)
}

func loadSportReadAllSetsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadAllSetsPath),
		contracts,
	)
}

func loadSportReadSetsByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadSetsByNamePath),
		contracts,
	)
}

func loadSportReadAllSetNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadAllSetNamesPath),
		contracts,
	)
}

//------------------------------------------------------------
// Plays
//------------------------------------------------------------
func loadSportCreatePlayTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportCreatePlayPath),
		contracts,
	)
}

func loadSportReadPlayByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadPlayByIDPath),
		contracts,
	)
}

func loadSportReadAllPlaysScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadAllPlaysPath),
		contracts,
	)
}

//------------------------------------------------------------
// Editions
//------------------------------------------------------------
func loadSportCreateEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportCreateEditionPath),
		contracts,
	)
}

func loadSportReadEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadEditionByIDPath),
		contracts,
	)
}

func loadSportCloseEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportCloseEditionPath),
		contracts,
	)
}

func loadSportReadAllEditionsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadAllEditionsPath),
		contracts,
	)
}

//------------------------------------------------------------
// Moment NFTs
//------------------------------------------------------------
func loadSportMintMomentNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportMintMomentNFTPath),
		contracts,
	)
}

func loadSportMintMomentNFTMultiTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportMintMomentNFTMultiPath),
		contracts,
	)
}

func loadSportReadMomentNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadMomentNFTSupplyPath),
		contracts,
	)
}

func loadSportReadMomentNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadMomentNFTPropertiesPath),
		contracts,
	)
}

func loadSportReadCollectionNFTLengthScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadCollectionNFTLengthPath),
		contracts,
	)
}

func loadSportReadCollectionNFTIDsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportReadCollectionNFTIDsPath),
		contracts,
	)
}

func loadSportTransferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportTransferNFTPath),
		contracts,
	)
}

func loadSportDisplayMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SportDisplayMetadataViewPath),
		contracts,
	)
}
