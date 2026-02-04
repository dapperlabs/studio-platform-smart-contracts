package atlas

import (
	_ "embed"
	"fmt"

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
	NFTStorefrontAddressV2 string
	ListingResourceIDs     []uint64
}

func (p DelistNFTStorefrontV2TxScriptParams) Validate() error {
	if p.NFTStorefrontAddressV2 == "" {
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

type ListNFTStorefrontV2Params struct {
	FungibleTokenContractAddress     string
	NonFungibleTokenContractAddress  string
	DapperUtilityCoinContractAddress string
	NFTProductName                   string
	NFTContractAddress               string
	NFTStorefrontV2ContractAddress   string
	NFTIDs                           string // comma-separated, e.g. "123,456"
	Prices                           string // comma-separated, e.g. "10.0,20.0"
}

func (p ListNFTStorefrontV2Params) Validate() error {
	if p.FungibleTokenContractAddress == "" ||
		p.NonFungibleTokenContractAddress == "" ||
		p.DapperUtilityCoinContractAddress == "" ||
		p.NFTProductName == "" ||
		p.NFTContractAddress == "" ||
		p.NFTStorefrontV2ContractAddress == "" ||
		p.NFTIDs == "" ||
		p.Prices == "" {
		return fmt.Errorf("all fields in AdminFulfillPackBuybackOfferParams must be non-empty")
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
