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
	DapperSportAddressPlaceholder   = "\"[^\"]*DapperSport.cdc\""
	metadataViewsAddressPlaceholder = "0xMETADATAVIEWSADDRESS"

	DapperSportPath                 = "../../../contracts/DapperSport.cdc"
	DapperSportTransactionsRootPath = "../../../transactions"
	DapperSportScriptsRootPath      = "../../../scripts"

	// Accounts
	DapperSportSetupAccountPath   = DapperSportTransactionsRootPath + "/user/setup_dappersport_account.cdc"
	DapperSportAccountIsSetupPath = DapperSportScriptsRootPath + "/user/account_is_setup.cdc"

	// Series
	DapperSportCreateSeriesPath       = DapperSportTransactionsRootPath + "/admin/series/create_series.cdc"
	DapperSportCloseSeriesPath        = DapperSportTransactionsRootPath + "/admin/series/close_series.cdc"
	DapperSportReadAllSeriesPath      = DapperSportScriptsRootPath + "/series/read_all_series.cdc"
	DapperSportReadSeriesByIDPath     = DapperSportScriptsRootPath + "/series/read_series_by_id.cdc"
	DapperSportReadSeriesByNamePath   = DapperSportScriptsRootPath + "/series/read_series_by_name.cdc"
	DapperSportReadAllSeriesNamesPath = DapperSportScriptsRootPath + "/series/read_all_series_names.cdc"

	// Sets
	DapperSportCreateSetPath       = DapperSportTransactionsRootPath + "/admin/sets/create_set.cdc"
	DapperSportLockSetPath       = DapperSportTransactionsRootPath + "/admin/sets/lock_set.cdc"
	DapperSportReadAllSetsPath     = DapperSportScriptsRootPath + "/sets/read_all_sets.cdc"
	DapperSportReadSetByIDPath     = DapperSportScriptsRootPath + "/sets/read_set_by_id.cdc"
	DapperSportReadSetsByNamePath  = DapperSportScriptsRootPath + "/sets/read_sets_by_name.cdc"
	DapperSportReadAllSetNamesPath = DapperSportScriptsRootPath + "/sets/read_all_set_names.cdc"

	// Plays
	DapperSportCreatePlayPath   = DapperSportTransactionsRootPath + "/admin/plays/create_play.cdc"
	DapperSportReadPlayByIDPath = DapperSportScriptsRootPath + "/plays/read_play_by_id.cdc"
	DapperSportReadAllPlaysPath = DapperSportScriptsRootPath + "/plays/read_all_plays.cdc"

	// Editions
	DapperSportCreateEditionPath   = DapperSportTransactionsRootPath + "/admin/editions/create_edition.cdc"
	DapperSportCloseEditionPath    = DapperSportTransactionsRootPath + "/admin/editions/close_edition.cdc"
	DapperSportReadEditionByIDPath = DapperSportScriptsRootPath + "/editions/read_edition_by_id.cdc"
	DapperSportReadAllEditionsPath = DapperSportScriptsRootPath + "/edition/read_all_editions.cdc"

	// Moment NFTs
	DapperSportMintMomentNFTPath                 = DapperSportTransactionsRootPath + "/admin/nfts/mint_moment_nft.cdc"
	DapperSportMintMomentNFTMultiPath            = DapperSportTransactionsRootPath + "/admin/nfts/mint_moment_nft_multi.cdc"
	DapperSportTransferNFTPath                   = DapperSportTransactionsRootPath + "/user/transfer_moment_nft.cdc"
	DapperSportReadMomentNFTSupplyPath           = DapperSportScriptsRootPath + "/nfts/read_moment_nft_supply.cdc"
	DapperSportReadMomentNFTPropertiesPath       = DapperSportScriptsRootPath + "/nfts/read_moment_nft_properties.cdc"
	DapperSportReadCollectionNFTLengthPath       = DapperSportScriptsRootPath + "/nfts/read_collection_nft_length.cdc"
	DapperSportReadCollectionNFTIDsPath          = DapperSportScriptsRootPath + "/nfts/read_collection_nft_ids.cdc"
	DapperSportDisplayMetadataViewPath           = DapperSportScriptsRootPath + "/nfts/metadata_display_view.cdc"
	DapperSportEditionMetadataViewPath           = DapperSportScriptsRootPath + "/nfts/metadata_editions_view.cdc"
	DapperSportSerialMetadataViewPath            = DapperSportScriptsRootPath + "/nfts/metadata_serial_view.cdc"
	DapperSportNFTCollectionDataMetadataViewPath = DapperSportScriptsRootPath + "/nfts/metadata_nft_collection_data_view.cdc"

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

	DapperSportRe := regexp.MustCompile(DapperSportAddressPlaceholder)
	code = DapperSportRe.ReplaceAll(code, []byte("0x"+contracts.DapperSportAddress.String()))

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

func LoadDapperSport(nftAddress flow.Address, metadataViewsAddr flow.Address) []byte {
	code := readFile(DapperSportPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddr.String()))

	return code
}

func loadDapperSportSetupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportSetupAccountPath),
		contracts,
	)
}

func loadDapperSportAccountIsSetupScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportAccountIsSetupPath),
		contracts,
	)
}

//------------------------------------------------------------
// Series
//------------------------------------------------------------
func loadDapperSportCreateSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportCreateSeriesPath),
		contracts,
	)
}

func loadDapperSportReadSeriesByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadSeriesByIDPath),
		contracts,
	)
}

func loadDapperSportReadSeriesByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadSeriesByNamePath),
		contracts,
	)
}

func loadDapperSportReadAllSeriesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadAllSeriesPath),
		contracts,
	)
}

func loadDapperSportReadAllSeriesNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadAllSeriesNamesPath),
		contracts,
	)
}

func loadDapperSportCloseSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportCloseSeriesPath),
		contracts,
	)
}

//------------------------------------------------------------
// Sets
//------------------------------------------------------------
func loadDapperSportCreateSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportCreateSetPath),
		contracts,
	)
}

func loadDapperSportLockSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportLockSetPath),
		contracts,
	)
}

func loadDapperSportReadSetByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadSetByIDPath),
		contracts,
	)
}

func loadDapperSportReadAllSetsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadAllSetsPath),
		contracts,
	)
}

func loadDapperSportReadSetsByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadSetsByNamePath),
		contracts,
	)
}

func loadDapperSportReadAllSetNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadAllSetNamesPath),
		contracts,
	)
}

//------------------------------------------------------------
// Plays
//------------------------------------------------------------
func loadDapperSportCreatePlayTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportCreatePlayPath),
		contracts,
	)
}

func loadDapperSportReadPlayByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadPlayByIDPath),
		contracts,
	)
}

func loadDapperSportReadAllPlaysScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadAllPlaysPath),
		contracts,
	)
}

//------------------------------------------------------------
// Editions
//------------------------------------------------------------
func loadDapperSportCreateEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportCreateEditionPath),
		contracts,
	)
}

func loadDapperSportReadEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadEditionByIDPath),
		contracts,
	)
}

func loadDapperSportCloseEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportCloseEditionPath),
		contracts,
	)
}

func loadDapperSportReadAllEditionsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadAllEditionsPath),
		contracts,
	)
}

//------------------------------------------------------------
// Moment NFTs
//------------------------------------------------------------
func loadDapperSportMintMomentNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportMintMomentNFTPath),
		contracts,
	)
}

func loadDapperSportMintMomentNFTMultiTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportMintMomentNFTMultiPath),
		contracts,
	)
}

func loadDapperSportReadMomentNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadMomentNFTSupplyPath),
		contracts,
	)
}

func loadDapperSportReadMomentNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadMomentNFTPropertiesPath),
		contracts,
	)
}

func loadDapperSportReadCollectionNFTLengthScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadCollectionNFTLengthPath),
		contracts,
	)
}

func loadDapperSportReadCollectionNFTIDsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportReadCollectionNFTIDsPath),
		contracts,
	)
}

func loadDapperSportTransferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportTransferNFTPath),
		contracts,
	)
}

func loadDapperSportDisplayMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportDisplayMetadataViewPath),
		contracts,
	)
}

func loadDapperSportEditionMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportEditionMetadataViewPath),
		contracts,
	)
}

func loadDapperSportSerialMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportSerialMetadataViewPath),
		contracts,
	)
}

func loadDapperSportNFTCollectionDataMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DapperSportNFTCollectionDataMetadataViewPath),
		contracts,
	)
}
