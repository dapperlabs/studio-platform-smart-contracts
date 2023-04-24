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
	DSSCollectionAddressPlaceholder = "\"[^\"]*DSSCollection.cdc\""
	metadataViewsAddressPlaceholder = "0xMETADATAVIEWSADDRESS"
	exampleNFTAddressPlaceholder    = "0xEXAMPLENFTADDRESS"
	addressUtilsPlaceholder         = "0xADDRESSUTILS"
)

const (
	DSSCollectionPath    = "../../../contracts/DSSCollection.cdc"
	ExampleNFTPath       = "../../../contracts/ExampleNFT.cdc"
	ArrayUtilsPath       = "../../../contracts/flow-utils/ArrayUtils.cdc"
	StringUtilsPath      = "../../../contracts/flow-utils/StringUtils.cdc"
	AddressUtilsPath     = "../../../contracts/flow-utils/AddressUtils.cdc"
	TransactionsRootPath = "../../../transactions"
	ScriptsRootPath      = "../../../scripts"

	// Accounts
	SetupAccountTxPath       = TransactionsRootPath + "/user/setup_user_account.cdc"
	IsAccountSetupScriptPath = ScriptsRootPath + "/user/is_account_setup.cdc"

	// Collection Groups
	CreateCollectionGroupTxPath          = TransactionsRootPath + "/admin/create_collection_group.cdc"
	CreateTimeBoundCollectionGroupTxPath = TransactionsRootPath + "/admin/create_collection_group_time_bound.cdc"
	CloseCollectionGroupTxPath           = TransactionsRootPath + "/admin/close_collection_group.cdc"
	GetCollectionGroupByIDScriptPath     = ScriptsRootPath + "/get_collection_group.cdc"
	GetCollectionGroupNFTCountScriptPath = ScriptsRootPath + "/get_collection_group_nft_count.cdc"
	CheckCollectionOwnershipScriptPath   = ScriptsRootPath + "/check_collection_ownership.cdc"

	// Slots
	CreateSlotTxPath       = TransactionsRootPath + "/admin/create_slot.cdc"
	GetSlotByIDScriptPath  = ScriptsRootPath + "/get_slot.cdc"
	CreateItemInSlotTxPath = TransactionsRootPath + "/admin/create_item_in_slot.cdc"

	// NFTs
	SetupExampleNFTxPath                 = TransactionsRootPath + "/user/setup_example_nft.cdc"
	TransferNFTTxPath                    = TransactionsRootPath + "/user/transfer_nft.cdc"
	MintNFTTxPath                        = TransactionsRootPath + "/admin/mint_nft.cdc"
	MintNFTTAndRecordxPath               = TransactionsRootPath + "/admin/mint_and_record.cdc"
	MintExampleNFTTxPath                 = TransactionsRootPath + "/admin/mint_example_nft.cdc"
	CompletedCollectionGroupTxPath       = TransactionsRootPath + "/admin/set_completed_collection_group.cdc"
	ReadNftSupplyScriptPath              = ScriptsRootPath + "/total_supply.cdc"
	ReadNftPropertiesTxPath              = ScriptsRootPath + "/get_nft.cdc"
	DSSCollectionDisplayMetadataViewPath = ScriptsRootPath + "/metadata_display_view.cdc"

	// MetadataViews
	MetadataViewsContractsBaseURL = "https://raw.githubusercontent.com/onflow/flow-nft/master/contracts/"
	MetadataViewsInterfaceFile    = "MetadataViews.cdc"
	MetadataFTReplaceAddress      = `"./utility/FungibleToken.cdc"`
	MetadataNFTReplaceAddress     = `"./NonFungibleToken.cdc"`

	// flow-utils
	StringUtilsAUReplaceAddress = `"./ArrayUtils.cdc"`
	StringUtilsSUReplaceAddress = `"./StringUtils.cdc"`
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func rX(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	DSSCollectionRe := regexp.MustCompile(DSSCollectionAddressPlaceholder)
	code = DSSCollectionRe.ReplaceAll(code, []byte("0x"+contracts.DSSCollectionAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewsAddress.String()))
	return code
}

func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	DapperSportRe := regexp.MustCompile(DSSCollectionAddressPlaceholder)
	code = DapperSportRe.ReplaceAll(code, []byte("0x"+contracts.DSSCollectionAddress.String()))

	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+contracts.MetadataViewsAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), exampleNFTAddressPlaceholder, "0x"+contracts.DSSCollectionAddress.String()))

	return code
}

func LoadDSSCollectionContract(nftAddress flow.Address, metadataViewsAddress flow.Address, addressUtilsAddr flow.Address) []byte {
	code := readFile(DSSCollectionPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), addressUtilsPlaceholder, "0x"+addressUtilsAddr.String()))
	return code
}

func LoadExampleNFTContract(nftAddress flow.Address, metadataViewsAddress flow.Address) []byte {
	code := readFile(ExampleNFTPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))
	code = []byte(strings.ReplaceAll(string(code), metadataViewsAddressPlaceholder, "0x"+metadataViewsAddress.String()))

	return code
}

func setupAccountTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SetupAccountTxPath),
		contracts,
	)
}

func setupExampleNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(SetupExampleNFTxPath),
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

func createSlotTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CreateSlotTxPath),
		contracts,
	)
}

func readCollectionGroupByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GetCollectionGroupByIDScriptPath),
		contracts,
	)
}

func readSlotByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GetSlotByIDScriptPath),
		contracts,
	)
}

func closeCollectionGroupTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CloseCollectionGroupTxPath),
		contracts,
	)
}

func createItemInSlotTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CreateItemInSlotTxPath),
		contracts,
	)
}

// ------------------------------------------------------------
// DSSCollection NFTs
// ------------------------------------------------------------
func transferNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(TransferNFTTxPath),
		contracts,
	)
}

func setCompletedCollectionGroup(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CompletedCollectionGroupTxPath),
		contracts,
	)
}

func mintDSSCollectionTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(MintNFTTxPath),
		contracts,
	)
}

func mintDSSCollectionAndRecordTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(MintNFTTAndRecordxPath),
		contracts,
	)
}

func mintExampleNFTTransaction(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(MintExampleNFTTxPath),
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

func LoadArrayUtils() []byte {
	return readFile(ArrayUtilsPath)
}

func LoadStringUtils(auAddress flow.Address) []byte {
	code := readFile(StringUtilsPath)
	code = []byte(strings.Replace(string(code), StringUtilsAUReplaceAddress, "0x"+auAddress.String(), 1))

	return code
}

func LoadAddressUtils(suAddress flow.Address) []byte {
	code := readFile(AddressUtilsPath)
	code = []byte(strings.Replace(string(code), StringUtilsSUReplaceAddress, "0x"+suAddress.String(), 1))

	return code
}

func loadDSSCollectionDisplayMetadataViewScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(DSSCollectionDisplayMetadataViewPath),
		contracts,
	)
}

func getCollectionGroupNFTCountScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GetCollectionGroupNFTCountScriptPath),
		contracts,
	)
}

func getCheckCollectionOwnershipScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(CheckCollectionOwnershipScriptPath),
		contracts,
	)
}
