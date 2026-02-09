package atlas

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/dapperlabs/studio-platform-smart-contracts/utils"
)

// Transactions is a list of all the transactions we export with imports mapped
var (
	//go:embed transactions/user/buy_packs_primary_sale.cdc
	UserBuyPacksPrimarySale []byte

	//go:embed transactions/user/fulfill_pack_buyback_offer.cdc
	AdminFulfillPackBuybackOffer []byte

	//go:embed transactions/user/delist_nftstorefront.cdc
	DelistNFTStorefront []byte

	//go:embed transactions/user/delist_nftstorefrontv2.cdc
	DelistNFTStorefrontV2 []byte

	//go:embed transactions/user/delist_topshotmarketv3.cdc
	DelistTopShotMarketV3 []byte

	//go:embed transactions/user/list_nftstorefrontv2.cdc
	ListNFTStorefrontV2 []byte
)

type UserBuyPacksPrimarySaleParams struct {
	FungibleTokenContractAddress     string
	NonFungibleTokenContractAddress  string
	DapperUtilityCoinContractAddress string
	NFTProductName                   string
	NFTContractAddress               string
	NFTStorefrontV1ContractAddress   string
	StorefrontAddress                string
	RecipientAddress                 string
	ListingIds                       string // comma-separated, e.g. "123,456"
}

func (p UserBuyPacksPrimarySaleParams) Validate() error {
	if p.FungibleTokenContractAddress == "" ||
		p.NonFungibleTokenContractAddress == "" ||
		p.DapperUtilityCoinContractAddress == "" ||
		p.NFTProductName == "" ||
		p.NFTContractAddress == "" ||
		p.NFTStorefrontV1ContractAddress == "" ||
		p.StorefrontAddress == "" ||
		p.RecipientAddress == "" ||
		p.ListingIds == "" {
		return fmt.Errorf("all fields in UserBuyPacksPrimarySaleParams must be non-empty")
	}
	return nil
}

func UserBuyPacksPrimarySaleTxScript(params UserBuyPacksPrimarySaleParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	bytes, err := utils.ParseCadenceTemplate(UserBuyPacksPrimarySale, params)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type AdminFulfillPackBuybackOfferParams struct {
	FungibleTokenContractAddress     string
	NonFungibleTokenContractAddress  string
	DapperUtilityCoinContractAddress string
	NFTProductName                   string
	NFTContractAddress               string
	NFTStorefrontV1ContractAddress   string
	NFTIDs                           string // comma-separated, e.g. "123,456"
	Prices                           string // comma-separated, e.g. "10.0,20.0"
}

func (p AdminFulfillPackBuybackOfferParams) Validate() error {
	if p.FungibleTokenContractAddress == "" ||
		p.NonFungibleTokenContractAddress == "" ||
		p.DapperUtilityCoinContractAddress == "" ||
		p.NFTProductName == "" ||
		p.NFTContractAddress == "" ||
		p.NFTStorefrontV1ContractAddress == "" ||
		p.NFTIDs == "" ||
		p.Prices == "" {
		return fmt.Errorf("all fields in AdminFulfillPackBuybackOfferParams must be non-empty")
	}
	return nil
}

func AdminFulfillPackBuybackOfferTxScript(params AdminFulfillPackBuybackOfferParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	bytes, err := utils.ParseCadenceTemplate(AdminFulfillPackBuybackOffer, params)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type DelistNFTStorefrontTxScriptParams struct {
	NFTStorefrontAddress string
	ListingResourceIDs   []uint64
}

func (p DelistNFTStorefrontTxScriptParams) Validate() error {
	if p.NFTStorefrontAddress == "" {
		return fmt.Errorf("NFTStorefrontAddress must be non-empty")
	}

	if len(p.ListingResourceIDs) == 0 {
		return fmt.Errorf("listingResourceIDs must contain at least one ID")
	}
	return nil
}

func DelistNFTStorefrontTxScript(params DelistNFTStorefrontTxScriptParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	bytes, err := utils.ParseCadenceTemplate(DelistNFTStorefront, params)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type DelistNFTStorefrontV2TxScriptParams struct {
	NFTStorefrontV2Address string
	ListingResourceIDs     []uint64
}

func (p DelistNFTStorefrontV2TxScriptParams) Validate() error {
	if p.NFTStorefrontV2Address == "" {
		return fmt.Errorf("NFTStorefrontAddress must be non-empty")
	}

	if len(p.ListingResourceIDs) == 0 {
		return fmt.Errorf("listingResourceIDs must contain at least one ID")
	}
	return nil
}

func DelistNFTStorefrontV2TxScript(params DelistNFTStorefrontV2TxScriptParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	bytes, err := utils.ParseCadenceTemplate(DelistNFTStorefrontV2, params)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type DelistTopShotMarketV3ScriptParams struct {
	TopShotContractAddress          string
	TopShotMarketContractAddress    string
	NonFungibleTokenContractAddress string
	NftIds                          []uint64
}

func (p DelistTopShotMarketV3ScriptParams) Validate() error {
	if p.TopShotContractAddress == "" ||
		p.TopShotMarketContractAddress == "" ||
		p.NonFungibleTokenContractAddress == "" ||
		len(p.NftIds) == 0 {
		return fmt.Errorf("all fields in DelistTopShotMarketV3ScriptParams must be non-empty")
	}
	return nil
}

func DelistTopShotMarketV3Script(params DelistTopShotMarketV3ScriptParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	bytes, err := utils.ParseCadenceTemplate(DelistTopShotMarketV3, params)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type ListNFTStorefrontV2Params struct {
	FungibleTokenContractAddress     string
	NonFungibleTokenContractAddress  string
	MetadataViewsAddress             string
	DapperUtilityCoinContractAddress string
	TokenForwardingContractAddress   string
	NFTProductName                   string
	NFTContractAddress               string
	NFTStorefrontV2ContractAddress   string
	NFTIDs                           []uint64
	PricesInCents                    []int
	SaleCommissionPercent            int    // e.g. 5 for 5%
	Expiry                           string // Unix timestamp as string
}

// NFTIDsString converts the NFTIDs from an array of uint64 to a string to be used by the template
func (p ListNFTStorefrontV2Params) NFTIDsString() string {
	nftIDStrings := make([]string, len(p.NFTIDs))
	for i, id := range p.NFTIDs {
		nftIDStrings[i] = fmt.Sprintf("%d", id)
	}
	return strings.Join(nftIDStrings, ", ")
}

// PricesString converts the PricesInCents from an array of int to a string to be used by the template
func (p ListNFTStorefrontV2Params) PricesString() string {
	floatPrices := make([]string, len(p.PricesInCents))
	for i, price := range p.PricesInCents {
		floatPrices[i] = centsToUFix64String(price)
	}
	return strings.Join(floatPrices, ",")
}

// SaleCommissionPercentString converts the SaleCommissionPercent from an int to a string to be used by the template
func (p ListNFTStorefrontV2Params) SaleCommissionPercentString() string {
	return centsToUFix64String(p.SaleCommissionPercent)
}

func (p ListNFTStorefrontV2Params) Validate() error {
	if p.FungibleTokenContractAddress == "" ||
		p.NonFungibleTokenContractAddress == "" ||
		p.MetadataViewsAddress == "" ||
		p.DapperUtilityCoinContractAddress == "" ||
		p.NFTProductName == "" ||
		p.NFTContractAddress == "" ||
		p.NFTStorefrontV2ContractAddress == "" ||
		len(p.NFTIDs) == 0 ||
		len(p.PricesInCents) == 0 {
		return fmt.Errorf("all fields in ListNFTStorefrontV2Params must be non-empty")
	}
	if len(p.NFTIDs) != len(p.PricesInCents) {
		return fmt.Errorf("NFTIDs and pricesInCents must have the same length")
	}
	return nil
}

func ListNFTStorefrontV2TxScript(params ListNFTStorefrontV2Params) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	bytes, err := utils.ParseCadenceTemplate(ListNFTStorefrontV2, params)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func centsToUFix64String(cents int) string {
	// cents -> units with 8 decimal places
	// 1 cent = 0.01 = 0.01000000
	return fmt.Sprintf("%d.%08d", cents/100, (cents%100)*1_000_000)
}
