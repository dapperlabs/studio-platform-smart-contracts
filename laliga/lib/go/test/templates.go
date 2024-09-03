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
	nftAddressPlaceholder              = "\"NonFungibleToken\""
	GolazosAddressPlaceholder          = "\"Golazos\""
	metadataViewsAddressPlaceholder    = "\"MetadataViews\""
	fungibleTokenAddressPlaceholder    = "\"FungibleToken\""
	nonFungibleTokenAddressPlaceholder = "\"NonFungibleToken\""
	burnerAddressPlaceholder           = "\"Burner\""
	viewResolverAddressPlaceholder     = "\"ViewResolver\""
	royaltyAddressPlaceholder          = "GOLAZOS_ROYALTY_ADDRESS"

	BurnerPath                  = "../../../contracts/imports/Burner.cdc"
	GolazosPath                 = "../../../contracts/Golazos.cdc"
	GolazosTransactionsRootPath = "../../../transactions"
	GolazosScriptsRootPath      = "../../../scripts"

	// Accounts
	GolazosSetupAccountPath   = GolazosTransactionsRootPath + "/user/setup_golazos_account.cdc"
	GolazosAccountIsSetupPath = GolazosScriptsRootPath + "/user/account_is_setup.cdc"

	// Series
	GolazosCreateSeriesPath       = GolazosTransactionsRootPath + "/admin/series/create_series.cdc"
	GolazosCloseSeriesPath        = GolazosTransactionsRootPath + "/admin/series/close_series.cdc"
	GolazosReadAllSeriesPath      = GolazosScriptsRootPath + "/series/read_all_series.cdc"
	GolazosReadSeriesByIDPath     = GolazosScriptsRootPath + "/series/read_series_by_id.cdc"
	GolazosReadSeriesByNamePath   = GolazosScriptsRootPath + "/series/read_series_by_name.cdc"
	GolazosReadAllSeriesNamesPath = GolazosScriptsRootPath + "/series/read_all_series_names.cdc"

	// Sets
	GolazosCreateSetPath       = GolazosTransactionsRootPath + "/admin/sets/create_set.cdc"
	GolazosLockSetPath         = GolazosTransactionsRootPath + "/admin/sets/lock_set.cdc"
	GolazosReadAllSetsPath     = GolazosScriptsRootPath + "/sets/read_all_sets.cdc"
	GolazosReadSetByIDPath     = GolazosScriptsRootPath + "/sets/read_set_by_id.cdc"
	GolazosReadSetsByNamePath  = GolazosScriptsRootPath + "/sets/read_sets_by_name.cdc"
	GolazosReadAllSetNamesPath = GolazosScriptsRootPath + "/sets/read_all_set_names.cdc"

	// Plays
	GolazosCreatePlayPath   = GolazosTransactionsRootPath + "/admin/plays/create_play.cdc"
	GolazosReadPlayByIDPath = GolazosScriptsRootPath + "/plays/read_play_by_id.cdc"
	GolazosReadAllPlaysPath = GolazosScriptsRootPath + "/plays/read_all_plays.cdc"

	// Editions
	GolazosCreateEditionPath   = GolazosTransactionsRootPath + "/admin/editions/create_edition.cdc"
	GolazosCloseEditionPath    = GolazosTransactionsRootPath + "/admin/editions/close_edition.cdc"
	GolazosReadEditionByIDPath = GolazosScriptsRootPath + "/editions/read_edition_by_id.cdc"
	GolazosReadAllEditionsPath = GolazosScriptsRootPath + "/edition/read_all_editions.cdc"

	// Moment NFTs
	GolazosMintMomentNFTPath                 = GolazosTransactionsRootPath + "/admin/nfts/mint_moment_nft.cdc"
	GolazosMintMomentNFTMultiPath            = GolazosTransactionsRootPath + "/admin/nfts/mint_moment_nft_multi.cdc"
	GolazosTransferNFTPath                   = GolazosTransactionsRootPath + "/user/transfer_moment_nft.cdc"
	GolazosReadMomentNFTSupplyPath           = GolazosScriptsRootPath + "/nfts/read_moment_nft_supply.cdc"
	GolazosReadMomentNFTPropertiesPath       = GolazosScriptsRootPath + "/nfts/read_moment_nft_properties.cdc"
	GolazosReadCollectionNFTLengthPath       = GolazosScriptsRootPath + "/nfts/read_collection_nft_length.cdc"
	GolazosReadCollectionNFTIDsPath          = GolazosScriptsRootPath + "/nfts/read_collection_nft_ids.cdc"
	GolazosDisplayMetadataViewPath           = GolazosScriptsRootPath + "/nfts/metadata_display_view.cdc"
	GolazosEditionMetadataViewPath           = GolazosScriptsRootPath + "/nfts/metadata_editions_view.cdc"
	GolazosSerialMetadataViewPath            = GolazosScriptsRootPath + "/nfts/metadata_serial_view.cdc"
	GolazosNFTCollectionDataMetadataViewPath = GolazosScriptsRootPath + "/nfts/metadata_nft_collection_data_view.cdc"
	GolazosTraitsMetadataViewPath            = GolazosScriptsRootPath + "/nfts/metadata_traits_view.cdc"

	// MetadataViews
	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"
	MetadataFTReplaceAddress      = `"./utility/FungibleToken.cdc"`
	MetadataNFTReplaceAddress     = `"./NonFungibleToken.cdc"`
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	GolazosRe := regexp.MustCompile(GolazosAddressPlaceholder)
	code = GolazosRe.ReplaceAll(code, []byte("0x"+contracts.GolazosAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), burnerAddressPlaceholder, "0x"+contracts.BurnerAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), viewResolverAddressPlaceholder, "0x"+contracts.ViewResolverAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), fungibleTokenAddressPlaceholder, "0x"+contracts.FtAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), nonFungibleTokenAddressPlaceholder, "0x"+contracts.NFTAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), royaltyAddressPlaceholder, "0x"+contracts.GolazosAddress.String()))

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

func LoadGolazos(nftAddress flow.Address, metadataViewsAddr flow.Address, ftAddress flow.Address, viewResolverAddress flow.Address) []byte {
	code := readFile(GolazosPath)

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddr.String()))
	code = []byte(strings.ReplaceAll(string(code), fungibleTokenAddressPlaceholder, "0x"+ftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), nonFungibleTokenAddressPlaceholder, "0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), viewResolverAddressPlaceholder, "0x"+viewResolverAddress.String()))

	return code
}

func loadGolazosSetupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosSetupAccountPath),
		contracts,
	)
}

func loadGolazosAccountIsSetupScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosAccountIsSetupPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Series
// ------------------------------------------------------------
func loadGolazosCreateSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosCreateSeriesPath),
		contracts,
	)
}

func loadGolazosReadSeriesByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadSeriesByIDPath),
		contracts,
	)
}

func loadGolazosReadSeriesByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadSeriesByNamePath),
		contracts,
	)
}

func loadGolazosReadAllSeriesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadAllSeriesPath),
		contracts,
	)
}

func loadGolazosReadAllSeriesNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadAllSeriesNamesPath),
		contracts,
	)
}

func loadGolazosCloseSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosCloseSeriesPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Sets
// ------------------------------------------------------------
func loadGolazosCreateSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosCreateSetPath),
		contracts,
	)
}

func loadGolazosLockSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosLockSetPath),
		contracts,
	)
}

func loadGolazosReadSetByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadSetByIDPath),
		contracts,
	)
}

func loadGolazosReadAllSetsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadAllSetsPath),
		contracts,
	)
}

func loadGolazosReadSetsByNameScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadSetsByNamePath),
		contracts,
	)
}

func loadGolazosReadAllSetNamesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadAllSetNamesPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Plays
// ------------------------------------------------------------
func loadGolazosCreatePlayTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosCreatePlayPath),
		contracts,
	)
}

func loadGolazosReadPlayByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadPlayByIDPath),
		contracts,
	)
}

func loadGolazosReadAllPlaysScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadAllPlaysPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Editions
// ------------------------------------------------------------
func loadGolazosCreateEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosCreateEditionPath),
		contracts,
	)
}

func loadGolazosReadEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadEditionByIDPath),
		contracts,
	)
}

func loadGolazosCloseEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosCloseEditionPath),
		contracts,
	)
}

func loadGolazosReadAllEditionsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadAllEditionsPath),
		contracts,
	)
}

// ------------------------------------------------------------
// Moment NFTs
// ------------------------------------------------------------
func loadGolazosMintMomentNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosMintMomentNFTPath),
		contracts,
	)
}

func loadGolazosMintMomentNFTMultiTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosMintMomentNFTMultiPath),
		contracts,
	)
}

func loadGolazosReadMomentNFTSupplyScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadMomentNFTSupplyPath),
		contracts,
	)
}

func loadGolazosReadMomentNFTPropertiesScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadMomentNFTPropertiesPath),
		contracts,
	)
}

func loadGolazosReadCollectionNFTLengthScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadCollectionNFTLengthPath),
		contracts,
	)
}

func loadGolazosReadCollectionNFTIDsScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosReadCollectionNFTIDsPath),
		contracts,
	)
}

func loadGolazosTransferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosTransferNFTPath),
		contracts,
	)
}

func loadGolazosDisplayMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosDisplayMetadataViewPath),
		contracts,
	)
}

func loadGolazosEditionMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosEditionMetadataViewPath),
		contracts,
	)
}

func loadGolazosSerialMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosSerialMetadataViewPath),
		contracts,
	)
}

func loadGolazosNFTCollectionDataMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosNFTCollectionDataMetadataViewPath),
		contracts,
	)
}

func loadGolazosTraitsMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GolazosTraitsMetadataViewPath),
		contracts,
	)
}
