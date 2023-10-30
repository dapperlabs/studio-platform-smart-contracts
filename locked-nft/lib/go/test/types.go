package test

import (
	"github.com/onflow/cadence"
)

type LockedData struct {
	Id          uint64
	LockedAt    uint64
	LockedUntil uint64
}

func parseLockedData(value cadence.Value) LockedData {
	fields := value.(cadence.Struct).Fields
	return LockedData{
		fields[0].ToGoValue().(uint64),
		fields[2].ToGoValue().(uint64),
		fields[3].ToGoValue().(uint64),
	}
}
