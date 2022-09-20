package test

import (
	"testing"
)

//------------------------------------------------------------
// Setup
//------------------------------------------------------------
func TestPDSDeployContracts(t *testing.T) {
	b := newEmulator()
	PDSDeployContracts(t, b)
}
