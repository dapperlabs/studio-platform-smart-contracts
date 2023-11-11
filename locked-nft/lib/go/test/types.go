package test

import (
	"github.com/onflow/cadence"
)

type LockedData struct {
	ID          uint64
	Owner       string
	LockedAt    uint64
	LockedUntil uint64
	nftType     string
}

func parseLockedData(value cadence.Value) LockedData {
	fields := value.(cadence.Optional).Value.(cadence.Struct).Fields

	return LockedData{
		fields[0].ToGoValue().(uint64),
		fields[1].(cadence.Address).String(),
		fields[2].ToGoValue().(uint64),
		fields[3].ToGoValue().(uint64),
		fields[5].(cadence.TypeValue).StaticType.ID(),
	}
}
