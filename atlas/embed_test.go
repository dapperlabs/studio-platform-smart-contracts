package atlas

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildFlowTxScript_AdminFulfillPackBuybackOffer(t *testing.T) {
	params := AdminFulfillPackBuybackOfferParams{
		FungibleTokenContractAddress:     "0xFUNGIBLE",
		NonFungibleTokenContractAddress:  "0xNONFUNGIBLE",
		DapperUtilityCoinContractAddress: "0xDUC",
		NFTProductName:                   "PackNFT",
		NFTContractAddress:               "0xNFT",
		NFTStorefrontV1ContractAddress:   "0xSTORE",
		NFTIDs:                           "789,101",
		Prices:                           "10.0,20.0",
	}

	tx, err := AdminFulfillPackBuybackOfferTxScript(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(tx, "0xFUNGIBLE") || !strings.Contains(tx, "10.0,20.0") {
		t.Errorf("expected substituted values in output, got: %s", tx)
	}
}

func TestBuildFlowTxScript_EmptyParamsError(t *testing.T) {
	params := AdminFulfillPackBuybackOfferParams{}
	tx, err := AdminFulfillPackBuybackOfferTxScript(params)
	if err == nil {
		t.Errorf("expected error for empty params, got nil (tx: %s)", tx)
	}
	if tx != "" {
		t.Errorf("expected empty tx string for error case, got: %s", tx)
	}
}

func TestDelistNFTStorefrontOutput(t *testing.T) {
	params := DelistNFTStorefrontTxScriptParams{
		NFTStorefrontAddress: "123",
		ListingResourceIDs:   []uint64{123, 456, 789},
	}

	tx, err := DelistNFTStorefrontTxScript(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("Generated transaction script:")
	fmt.Println(tx)
	t.Logf("Generated output:\n%s", tx)
}

func TestDelistTopShotMarketV3Output(t *testing.T) {
	params := DelistTopShotMarketV3ScriptParams{
		TopShotContractAddress:          "TOPSHOT",
		TopShotMarketContractAddress:    "MARKET",
		NonFungibleTokenContractAddress: "NFT",
		NftIds:                          []uint64{123, 456, 789},
	}

	tx, err := DelistTopShotMarketV3Script(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("Generated transaction script:")
	fmt.Println(tx)
	t.Logf("Generated output:\n%s", tx)
}

func TestDelistNFTStorefrontV2Output(t *testing.T) {
	// Test what Go's text/template outputs when given a []uint64
	params := DelistNFTStorefrontV2TxScriptParams{
		NFTStorefrontV2Address: "123",
		ListingResourceIDs:     []uint64{123, 456, 789},
	}

	tx, err := DelistNFTStorefrontV2TxScript(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("Generated transaction script:")
	fmt.Println(tx)
	t.Logf("Generated output:\n%s", tx)
}
