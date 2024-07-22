package test

import (
	"regexp"

	"github.com/onflow/flow-go-sdk"
)

// Handle relative paths by making these regular expressions

const (
	nftProviderAggregatorAddressPlaceholder = "\"NFTProviderAggregator\""
	exampleNFTAddressPlaceholder            = "\"ExampleNFT\""
	nftAddressPlaceholder                   = "\"NonFungibleToken\""
	metadataViewsAddressPlaceholder         = "\"MetadataViews\""
	viewResolverAddressPlaceholder          = "\"ViewResolver\""
	ftAddressPlaceholder                    = "\"FungibleToken\""
	burnerPlaceholder                       = "\"Burner\""

	NFTProviderAggregatorPath = "../../../contracts/NFTProviderAggregator.cdc"
	BurnerContractPath        = "../../../contracts/imports/Burner.cdc"
	TransactionsRootPath      = "../../../transactions"
	ScriptsRootPath           = "../../../scripts"
	NonFungibleTokenPath      = "../../../contracts/imports/NonFungibleToken.cdc"

	// ExampleNFT Transactions
	ExampleNftLinkWithdrawCapPath   = TransactionsRootPath + "/exampleNFT/link_providerCap_exampleNFT.cdc"
	ExampleNftMintPath              = TransactionsRootPath + "/exampleNFT/mint_exampleNFT.cdc"
	ExampleNftMintBatchedPath       = TransactionsRootPath + "/exampleNFT/mint_exampleNFTBatched.cdc"
	ExampleNftSetupPath             = TransactionsRootPath + "/exampleNFT/setup_exampleNFT.cdc"
	ExampleNftTransferPath          = TransactionsRootPath + "/exampleNFT/transfer_exampleNFT.cdc"
	ExampleNftRevokeWithdrawCapPath = TransactionsRootPath + "/exampleNFT/revoke_withdraw_cap.cdc"

	// ExampleNFT Scripts
	ExampleNftGetIdsPath = ScriptsRootPath + "/exampleNFT/balance_exampleNFT.cdc"

	// NFTProviderAggregator Transactions
	AddNftWithdrawCapAsManagerPath                    = TransactionsRootPath + "/add_nft_provider_capability_as_manager.cdc"
	AddNftWithdrawCapAsSupplierPath                   = TransactionsRootPath + "/add_nft_provider_capability_as_supplier.cdc"
	BootstrapAggregatorPath                           = TransactionsRootPath + "/bootstrap_aggregator.cdc"
	BootstrapSupplierPath                             = TransactionsRootPath + "/bootstrap_supplier.cdc"
	ClaimAggregatedNftWithdrawCapPath                 = TransactionsRootPath + "/claim_aggregated_nft_provider_capability.cdc"
	DestroyAggregatorPath                             = TransactionsRootPath + "/destroy_aggregator.cdc"
	DestroySupplierPath                               = TransactionsRootPath + "/destroy_supplier.cdc"
	PublishAdditionalSupplierFactoryCapPath           = TransactionsRootPath + "/publish_additional_supplier_factory_capabilities.cdc"
	PublishAggregatedNftWithdrawCapPath               = TransactionsRootPath + "/publish_aggregated_nft_provider_capability.cdc"
	RemoveNftWithdrawCapAsManagerPath                 = TransactionsRootPath + "/remove_nft_provider_capability_as_manager.cdc"
	RemoveNftWithdrawCapAsSupplierPath                = TransactionsRootPath + "/remove_nft_provider_capability_as_supplier.cdc"
	TransferFromAggregatedNftProviderAsManagerPath    = TransactionsRootPath + "/transfer_from_aggregated_nft_provider_as_manager.cdc"
	TransferFromAggregatedNftProviderAsThirdPartyPath = TransactionsRootPath + "/transfer_from_aggregated_nft_provider_as_thirdparty.cdc"
	UnpublishInboxCapPath                             = TransactionsRootPath + "/unpublish_inbox_capability.cdc"

	// NFTProviderAggregator Scripts
	GetAggregatorUuidPath               = ScriptsRootPath + "/get_aggregator_uuid.cdc"
	GetCollectionUuidsPath              = ScriptsRootPath + "/get_collection_uuids.cdc"
	GetManagerCollectionUuidsPath       = ScriptsRootPath + "/get_manager_collection_uuids.cdc"
	GetIdsPath                          = ScriptsRootPath + "/get_ids.cdc"
	GetSupplierAddedCollectionUuidsPath = ScriptsRootPath + "/get_supplier_added_collection_uuids.cdc"
)

// ------------------------------------------------------------
// Accounts
// ------------------------------------------------------------
func replaceAddresses(code []byte, contracts Contracts) []byte {
	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+contracts.NFTAddress.String()))

	nftProviderAggregatorRe := regexp.MustCompile(nftProviderAggregatorAddressPlaceholder)
	code = nftProviderAggregatorRe.ReplaceAll(code, []byte("0x"+contracts.NFTProviderAggregatorAddress.String()))

	exampleNftRe := regexp.MustCompile(exampleNFTAddressPlaceholder)
	code = exampleNftRe.ReplaceAll(code, []byte("0x"+contracts.ExampleNFTAddress.String()))

	metadataViewsRe := regexp.MustCompile(metadataViewsAddressPlaceholder)
	code = metadataViewsRe.ReplaceAll(code, []byte("0x"+contracts.MetadataViewAddress.String()))

	viewResolverRe := regexp.MustCompile(viewResolverAddressPlaceholder)
	code = viewResolverRe.ReplaceAll(code, []byte("0x"+contracts.ViewResolverAddress.String()))

	ftRe := regexp.MustCompile(ftAddressPlaceholder)
	code = ftRe.ReplaceAll(code, []byte("0x"+ftAddress.String()))

	burnerRe := regexp.MustCompile(burnerPlaceholder)
	code = burnerRe.ReplaceAll(code, []byte("0x"+contracts.BurnerAddress.String()))

	return code
}

func LoadNftProviderAggregator(nftAddress, metadataViewsAddr, viewResolverAddr, burnerAddr flow.Address) []byte {
	code := readFile(NFTProviderAggregatorPath)

	nftRe := regexp.MustCompile(nftAddressPlaceholder)
	code = nftRe.ReplaceAll(code, []byte("0x"+nftAddress.String()))

	metadataViewsRe := regexp.MustCompile(metadataViewsAddressPlaceholder)
	code = metadataViewsRe.ReplaceAll(code, []byte("0x"+metadataViewsAddr.String()))

	viewResolverRe := regexp.MustCompile(viewResolverAddressPlaceholder)
	code = viewResolverRe.ReplaceAll(code, []byte("0x"+viewResolverAddr.String()))

	ftRe := regexp.MustCompile(ftAddressPlaceholder)
	code = ftRe.ReplaceAll(code, []byte("0x"+ftAddress.String()))

	burnerRe := regexp.MustCompile(burnerPlaceholder)
	code = burnerRe.ReplaceAll(code, []byte("0x"+burnerAddr.String()))

	return code
}

func loadScript(contracts Contracts, path string) []byte {
	return replaceAddresses(
		readFile(path),
		contracts,
	)
}
