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
