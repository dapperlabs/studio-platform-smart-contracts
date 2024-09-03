package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk"
)

func getExampleNFTCollectionIds(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	address flow.Address,
) []interface{} {
	script := loadScript(contracts, ExampleNftGetIdsPath)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.BytesToAddress(address.Bytes()))})

	return GetFieldValue(result).([]interface{})
}

func getCollectionNftIds(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	address flow.Address,
) []interface{} {
	script := loadScript(contracts, GetIdsPath)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.BytesToAddress(address.Bytes()))})

	return GetFieldValue(result).([]interface{})
}

func getManagerCollectionUUIDs(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	address flow.Address,
) []interface{} {
	script := loadScript(contracts, GetManagerCollectionUuidsPath)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.BytesToAddress(address.Bytes()))})

	return GetFieldValue(result).([]interface{})
}
