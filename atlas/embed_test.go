package atlas

import (
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

func TestBuildFlowTxScript_Generic(t *testing.T) {
	tmpl := []byte("Hello, {{.Name}}!")
	params := struct{ Name string }{Name: "World"}
	result, err := BuildFlowTxScript(tmpl, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("unexpected template output: %s", result)
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
