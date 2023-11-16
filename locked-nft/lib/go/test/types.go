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

func parseInventoryData(value cadence.Value) []uint64 {
	values := value.(cadence.Optional).Value.(cadence.Array).Values

	inventory := make([]uint64, len(values))
	for i, v := range values {
		inventory[i] = uint64(v.(cadence.UInt64))
	}
	return inventory
}
