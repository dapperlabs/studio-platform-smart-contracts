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

func readLockedTokenByIDScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GetLockedTokenByIDScriptPath),
		contracts,
	)
}

func readInventoryScript(contracts Contracts) []byte {
	return replaceAddresses(
		readFile(GetInventoryScriptPath),
		contracts,
	)
}

func getLockedTokenData(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	id uint64,
) LockedData {
	script := readLockedTokenByIDScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.UInt64(id))})

	return parseLockedData(result)
}

func getNFTInventory(
	t *testing.T,
	b *emulator.Blockchain,
	contracts Contracts,
	userAddress flow.Address,
) []uint64 {
	script := readInventoryScript(contracts)
	result := executeScriptAndCheck(t, b, script, [][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))})
	var nftIds []uint64
	for _, val := range result.(cadence.Array).Values {
		nftId := val.ToGoValue().(uint64)
		nftIds = append(nftIds, nftId)
	}

	return nftIds
}
