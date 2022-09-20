package test

import (
	"testing"
)

//------------------------------------------------------------
// Setup
//------------------------------------------------------------
func TestPackNftDeployContracts(t *testing.T) {
	b := newEmulator()
	PackNftDeployContracts(t, b)
}
