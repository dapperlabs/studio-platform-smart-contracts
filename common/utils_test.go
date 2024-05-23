package common

import (
	"github.com/dapperlabs/studio-platform-smart-contracts/pds"
	"testing"
)

func TestGetUserOwnedEditionToMomentMap_Transaction(t *testing.T) {
	script, err := pds.Transaction.ReadFile("transactions/deploy/deploy-packNFT-with-auth.cdc")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(script))
}

func TestGetUserOwnedEditionToMomentMap_Script(t *testing.T) {
	script, err := pds.Scripts.ReadFile("scripts/packNFT/balance_packNFT.cdc")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(script))
	dddd := NewReplacer()
	aaa, err := dddd.ReplaceAddresses(string(script), struct {
		NonFungibleToken string
	}{
		"0x1",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(aaa)
}
