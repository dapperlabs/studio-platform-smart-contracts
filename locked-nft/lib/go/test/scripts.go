package test

import (
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-emulator/emulator"
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

	return GetFieldValue(result).(bool)
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
