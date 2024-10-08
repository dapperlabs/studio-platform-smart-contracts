package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk/test"

	"github.com/onflow/flow-nft/lib/go/contracts"
)

const addrA = "0x0A"

func TestNonFungibleTokenContract(t *testing.T) {
	addresses := test.AddressGenerator()
	addressA := addresses.New()
	contract := contracts.NonFungibleToken(addressA.String())
	assert.NotNil(t, contract)
	assert.Contains(t, string(contract), addressA.String())
}

func TestExampleNFTContract(t *testing.T) {
	addresses := test.AddressGenerator()
	addressA := addresses.New()
	addressB := addresses.New()
	addressC := addresses.New()

	contract := contracts.ExampleNFT(addressA, addressB, addressC)
	assert.NotNil(t, contract)

	assert.Contains(t, string(contract), addressA.String())
	assert.Contains(t, string(contract), addressB.String())
	assert.Contains(t, string(contract), addressC.String())
}

func TestMetadataViewsContract(t *testing.T) {
	addresses := test.AddressGenerator()
	addressA := addresses.New()
	addressB := addresses.New()
	addressC := addresses.New()

	contract := contracts.MetadataViews(addressA.String(), addressB.String(), addressC.String())
	assert.NotNil(t, contract)
}
