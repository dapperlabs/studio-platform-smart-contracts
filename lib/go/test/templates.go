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
	nftAddressPlaceholder    = "\"[^\"]*NonFungibleToken.cdc\""
	GolazoAddressPlaceholder = "\"[^\"]*Golazo.cdc\""
	metadataViewsAddressPlaceholder = "METADATAVIEWSADDRESS"

	GolazoPath                 = "../../../contracts/Golazo.cdc"
	GolazoTransactionsRootPath = "../../../transactions"
	GolazoScriptsRootPath      = "../../../scripts"

	// Accounts
	GolazoSetupAccountPath   = GolazoTransactionsRootPath + "/user/setup_Golazo_account.cdc"
	GolazoAccountIsSetupPath = GolazoScriptsRootPath + "/user/account_is_setup.cdc"

	// Series
	GolazoCreateSeriesPath       = GolazoTransactionsRootPath + "/admin/series/create_series.cdc"
	GolazoCloseSeriesPath        = GolazoTransactionsRootPath + "/admin/series/close_series.cdc"
	GolazoReadAllSeriesPath      = GolazoScriptsRootPath + "/series/read_all_series.cdc"
	GolazoReadSeriesByIDPath     = GolazoScriptsRootPath + "/series/read_series_by_id.cdc"
	GolazoReadSeriesByNamePath   = GolazoScriptsRootPath + "/series/read_series_by_name.cdc"
	GolazoReadAllSeriesNamesPath = GolazoScriptsRootPath + "/series/read_all_series_names.cdc"

	// Sets
	GolazoCreateSetPath       = GolazoTransactionsRootPath + "/admin/sets/create_set.cdc"
	GolazoReadAllSetsPath     = GolazoScriptsRootPath + "/sets/read_all_sets.cdc"
	GolazoReadSetByIDPath     = GolazoScriptsRootPath + "/sets/read_set_by_id.cdc"
	GolazoReadSetsByNamePath  = GolazoScriptsRootPath + "/sets/read_sets_by_name.cdc"
	GolazoReadAllSetNamesPath = GolazoScriptsRootPath + "/sets/read_all_set_names.cdc"

	// Plays
	GolazoCreatePlayPath   = GolazoTransactionsRootPath + "/admin/plays/create_play.cdc"
	GolazoReadPlayByIDPath = GolazoScriptsRootPath + "/plays/read_play_by_id.cdc"
	GolazoReadAllPlaysPath = GolazoScriptsRootPath + "/plays/read_all_plays.cdc"

	// Editions
	GolazoCreateEditionPath   = GolazoTransactionsRootPath + "/admin/editions/create_edition.cdc"
	GolazoCloseEditionPath    = GolazoTransactionsRootPath + "/admin/editions/close_edition.cdc"
	GolazoReadEditionByIDPath = GolazoScriptsRootPath + "/editions/read_edition_by_id.cdc"
	GolazoReadAllEditionsPath = GolazoScriptsRootPath + "/edition/read_all_editions.cdc"

	// Moment NFTs
	GolazoMintMomentNFTPath           = GolazoTransactionsRootPath + "/admin/nfts/mint_moment_nft.cdc"
	GolazoMintMomentNFTMultiPath      = GolazoTransactionsRootPath + "/admin/nfts/mint_moment_nft_multi.cdc"
	GolazoTransferNFTPath             = GolazoTransactionsRootPath + "/user/transfer_moment_nft.cdc"
	GolazoReadMomentNFTSupplyPath     = GolazoScriptsRootPath + "/nfts/read_moment_nft_supply.cdc"
	GolazoReadMomentNFTPropertiesPath = GolazoScriptsRootPath + "/nfts/read_moment_nft_properties.cdc"
	GolazoReadCollectionNFTLengthPath = GolazoScriptsRootPath + "/nfts/read_collection_nft_length.cdc"
	GolazoReadCollectionNFTIDsPath    = GolazoScriptsRootPath + "/nfts/read_collection_nft_ids.cdc"

	// MetadataViews
	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"
	MetadataFTReplaceAddress  = `"./utility/FungibleToken.cdc"`
	MetadataNFTReplaceAddress = `"./NonFungibleToken.cdc"`
)

//------------------------------------------------------------
// Accounts
//------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	GolazoRe := regexp.MustCompile(GolazoAddressPlaceholder)
	code = GolazoRe.ReplaceAll(code, []byte("0x"+contracts.GolazoAddress.String()))

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
	code = []byte(strings.Replace(strings.Replace(string(code), MetadataFTReplaceAddress, "0x"+ftAddress.String(), 1),MetadataNFTReplaceAddress, "0x"+nftAddress.String(), 1))

	return code
}

func LoadGolazo(nftAddress flow.Address, metadataViewsAddr flow.Address) []byte {
	code := readFile(GolazoPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, metadataViewsAddr.String()))

	return code
}

func loadGolazoSetupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoSetupAccountPath),
		contracts,
	)
}

func loadGolazoAccountIsSetupScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoAccountIsSetupPath),
		contracts,
	)
}

//------------------------------------------------------------
// Series
//------------------------------------------------------------
func loadGolazoCreateSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoCreateSeriesPath),
		contracts,
	)
}

func loadGolazoReadSeriesByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadSeriesByIDPath),
		contracts,
	)
}

func loadGolazoReadSeriesByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadSeriesByNamePath),
		contracts,
	)
}

func loadGolazoReadAllSeriesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadAllSeriesPath),
		contracts,
	)
}

func loadGolazoReadAllSeriesNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadAllSeriesNamesPath),
		contracts,
	)
}

func loadGolazoCloseSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoCloseSeriesPath),
		contracts,
	)
}

//------------------------------------------------------------
// Sets
//------------------------------------------------------------
func loadGolazoCreateSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoCreateSetPath),
		contracts,
	)
}

func loadGolazoReadSetByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadSetByIDPath),
		contracts,
	)
}

func loadGolazoReadAllSetsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadAllSetsPath),
		contracts,
	)
}

func loadGolazoReadSetsByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadSetsByNamePath),
		contracts,
	)
}

func loadGolazoReadAllSetNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadAllSetNamesPath),
		contracts,
	)
}

//------------------------------------------------------------
// Plays
//------------------------------------------------------------
func loadGolazoCreatePlayTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoCreatePlayPath),
		contracts,
	)
}

func loadGolazoReadPlayByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadPlayByIDPath),
		contracts,
	)
}

func loadGolazoReadAllPlaysScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadAllPlaysPath),
		contracts,
	)
}

//------------------------------------------------------------
// Editions
//------------------------------------------------------------
func loadGolazoCreateEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoCreateEditionPath),
		contracts,
	)
}

func loadGolazoReadEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadEditionByIDPath),
		contracts,
	)
}

func loadGolazoCloseEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoCloseEditionPath),
		contracts,
	)
}

func loadGolazoReadAllEditionsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadAllEditionsPath),
		contracts,
	)
}

//------------------------------------------------------------
// Moment NFTs
//------------------------------------------------------------
func loadGolazoMintMomentNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoMintMomentNFTPath),
		contracts,
	)
}

func loadGolazoMintMomentNFTMultiTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoMintMomentNFTMultiPath),
		contracts,
	)
}

func loadGolazoReadMomentNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadMomentNFTSupplyPath),
		contracts,
	)
}

func loadGolazoReadMomentNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadMomentNFTPropertiesPath),
		contracts,
	)
}

func loadGolazoReadCollectionNFTLengthScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadCollectionNFTLengthPath),
		contracts,
	)
}

func loadGolazoReadCollectionNFTIDsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoReadCollectionNFTIDsPath),
		contracts,
	)
}

func loadGolazoTransferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazoTransferNFTPath),
		contracts,
	)
}
