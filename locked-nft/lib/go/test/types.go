package test

import (
	"github.com/onflow/cadence"
)

type LockedData struct {
	Owner       string
	LockedAt    uint64
	LockedUntil uint64
	nftType     string
}

func parseLockedData(value cadence.Value) LockedData {
	fields := value.(cadence.Struct).Fields
	return LockedData{
		fields[0].ToGoValue().(string),
		fields[1].ToGoValue().(uint64),
		fields[2].ToGoValue().(uint64),
		fields[3].ToGoValue().(string),
	}
}
