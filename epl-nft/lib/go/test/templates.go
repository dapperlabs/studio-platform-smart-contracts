package test

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/onflow/flow-go-sdk"
)

const (
	nftAddressPlaceholder           = "\"[^\"]*NonFungibleToken.cdc\""
	EPLAddressPlaceholder           = "\"[^\"]*EnglishPremierLeague.cdc\""
	metadataViewsAddressPlaceholder = "0xMETADATAVIEWSADDRESS"
	ftAddressPlaceholder            = "\"[^\"]*FungibleToken.cdc\""
)

const (
	EPLPath              = "../../../contracts/EnglishPremierLeague.cdc"
	TransactionsRootPath = "../../../transactions"
	ScriptsRootPath      = "../../../scripts"

	// Accounts
	SetupAccountTxPath       = TransactionsRootPath + "/user/setup_collection.cdc"
	IsAccountSetupScriptPath = ScriptsRootPath + "/is_account_setup.cdc"

	// NFTs
	MintNFTTxPath              = TransactionsRootPath + "/admin/mint_nft.cdc"
	ReadNftSupplyScriptPath    = ScriptsRootPath + "/total_supply.cdc"
	ReadNftPropertiesTxPath    = ScriptsRootPath + "/get_nft.cdc"
	EPLDisplayMetadataViewPath = ScriptsRootPath + "/metadata_display_view.cdc"

	// MetadataViews
	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"
	MetadataFTReplaceAddress      = `"./utility/FungibleToken.cdc"`
	MetadataNFTReplaceAddress     = `"./NonFungibleToken.cdc"`

	// FungibleToken
	FTPath = "../../../contracts/FungibleToken.cdc"

	// Series
	EPLCreateSeriesPath   = TransactionsRootPath + "/admin/create_series.cdc"
	EPLCloseSeriesPath    = TransactionsRootPath + "/admin/close_series.cdc"
	EPLReadSeriesByIDPath = ScriptsRootPath + "/get_series.cdc"

	// Sets
	EPLCreateSetPath   = TransactionsRootPath + "/admin/create_set.cdc"
	EPLLockSetPath     = TransactionsRootPath + "/admin/lock_set.cdc"
	EPLReadSetByIDPath = ScriptsRootPath + "/get_set.cdc"

	// Tag
	EPLCreateTagPath   = TransactionsRootPath + "/admin/create_tag.cdc"
	EPLReadTagByIDPath = ScriptsRootPath + "/get_tag.cdc"

	// Plays
	EPLCreatePlayPath   = TransactionsRootPath + "/admin/create_play.cdc"
	EPLReadPlayByIDPath = ScriptsRootPath + "/get_play.cdc"

	// Editions
	EPLCreateEditionPath   = TransactionsRootPath + "/admin/create_edition.cdc"
	EPLCloseEditionPath    = TransactionsRootPath + "/admin/close_edition.cdc"
	EPLReadEditionByIDPath = ScriptsRootPath + "/get_edition.cdc"
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func rX(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	DSSCollectionRe := regexp.MustCompile(EPLAddressPlaceholder)
	code = DSSCollectionRe.ReplaceAll(code, []byte("0x"+contracts.EPLAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewsAddress.String()))

	return code
}

func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	DapperSportRe := regexp.MustCompile(EPLAddressPlaceholder)
	code = DapperSportRe.ReplaceAll(code, []byte("0x"+contracts.EPLAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewsAddress.String()))

	return code
}

func LoadEPLContract(nftAddress flow.Address, metadataViewsAddress flow.Address) []byte {
	code := readFile(EPLPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	ftRe := regexp.MustCompile(ftAddressPlaceholder)
	code = ftRe.ReplaceAll(code, []byte("0x"+ftAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddress.String()))

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

func loadEPLReadSeriesByIdScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLReadSeriesByIDPath),
		contracts,
	)
}

func loadEPLCreateSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCreateSeriesPath),
		contracts,
	)
}

func loadEPLCloseSeriesTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCloseSeriesPath),
		contracts,
	)
}

func loadEPLCreateSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCreateSetPath),
		contracts,
	)
}

func loadEPLLockSetTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLLockSetPath),
		contracts,
	)
}

func loadEPLReadSetByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLReadSetByIDPath),
		contracts,
	)
}

func loadEPLReadTagByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLReadTagByIDPath),
		contracts,
	)
}

func loadEPLCreateTagTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCreateTagPath),
		contracts,
	)
}

func loadEPLCreatePlayTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCreatePlayPath),
		contracts,
	)
}

func loadEPLReadPlayByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLReadPlayByIDPath),
		contracts,
	)
}

func loadEPLCreateEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCreateEditionPath),
		contracts,
	)
}

func loadEPLReadEditionByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLReadEditionByIDPath),
		contracts,
	)
}

func loadEPLCloseEditionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLCloseEditionPath),
		contracts,
	)
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

func loadEPLDisplayMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(EPLDisplayMetadataViewPath),
		contracts,
	)
}
