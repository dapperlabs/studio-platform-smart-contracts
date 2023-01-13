package test

import (
	"github.com/onflow/cadence"
)

type CollectionGroupData struct {
	ID                     uint64
	Name                   string
	Open                   bool
	TimeBound              bool
	NFTIDInCollectionGroup map[uint64]bool
}

type NFTData struct {
	ID                uint64
	CollectionGroupID uint64
}

func parseCollectionGroupData(value cadence.Value) CollectionGroupData {
	fields := value.(cadence.Struct).Fields
	return CollectionGroupData{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(string),
		fields[3].ToGoValue().(bool),
		fields[6].ToGoValue().(bool),
		cadenceUint64DictToGo(fields[7].(cadence.Dictionary)),
	}
}

func parseDSSCollectionProperties(value cadence.Value) NFTData {
	array := value.(cadence.Array).Values
	return NFTData{
		array[0].ToGoValue().(uint64),
		array[1].ToGoValue().(uint64),
	}
}

func cadenceStringDictToGo(cadenceDict cadence.Dictionary) map[string]string {
	goDict := make(map[string]string)
	for _, pair := range cadenceDict.Pairs {
		goDict[pair.Key.ToGoValue().(string)] = pair.Value.ToGoValue().(string)
	}
	return goDict
}

func cadenceUint64DictToGo(cadenceDict cadence.Dictionary) map[uint64]bool {
	goDict := make(map[uint64]bool)
	for _, pair := range cadenceDict.Pairs {
		goDict[pair.Key.ToGoValue().(uint64)] = pair.Value.ToGoValue().(bool)
	}
	return goDict
}
