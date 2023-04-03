package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
)

func isAccountSetup(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	address flow.Address,
) bool {
	script := isAccountSetupScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.BytesToAddress(address.Bytes()))})

	return result.ToGoValue().(bool)
}

func getDSSCollectionNFTDisplayMetadataView(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	collectionAddress flow.Address,
	nftID uint64,
) DisplayView {
	script := loadDSSCollectionDisplayMetadataViewScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{
		jsoncdc.MustEncode(cadence.BytesToAddress(collectionAddress.Bytes())),
		jsoncdc.MustEncode(cadence.UInt64(nftID)),
	})

	return parseMetadataDisplayView(result)
}

func getCollectionGroupNFTCount(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	collectionGroupID uint64,
) uint64 {
	script := getCollectionGroupNFTCountScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{
		jsoncdc.MustEncode(cadence.UInt64(collectionGroupID)),
	})

	nftCount := result.ToGoValue().(uint64)
	return nftCount
}

func checkCollectionOwnership(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	isAuthorized bool,
	userAddress string,
	collectionGroupID uint64,
) bool {
	script := getCheckCollectionOwnershipScript(contracts)
	address := flow.HexToAddress(userAddress)
	result := executeScriptAndCheck(t, b, script, [][]byte{
		jsoncdc.MustEncode(cadence.Address(address)),
		jsoncdc.MustEncode(cadence.UInt64(collectionGroupID)),
	})

	return result.ToGoValue().(bool)
}
