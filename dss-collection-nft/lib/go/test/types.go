package test

import (
	"github.com/onflow/cadence"
)

type CollectionGroupData struct {
	ID          uint64
	Name        string
	Description string
	TypeName    string
	Open        bool
	TimeBound   bool
}

type NFTData struct {
	ID                uint64
	CollectionGroupID uint64
	SerialNumber      uint64
	CompletionDate    uint64
	CompletedBy       string
}

type DisplayView struct {
	Name        string
	Description string
	ImageURL    string
}

func parseCollectionGroupData(value cadence.Value) CollectionGroupData {
	fields := value.(cadence.Struct).Fields
	return CollectionGroupData{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(string),
		fields[2].ToGoValue().(string),
		fields[3].ToGoValue().(string),
		fields[4].ToGoValue().(bool),
		fields[7].ToGoValue().(bool),
	}
}

func parseNFTProperties(value cadence.Value) NFTData {
	array := value.(cadence.Array).Values
	return NFTData{
		array[0].ToGoValue().(uint64),
		array[1].ToGoValue().(uint64),
		array[2].ToGoValue().(uint64),
		array[3].ToGoValue().(uint64),
		array[4].ToGoValue().(string),
	}
}

func parseMetadataDisplayView(value cadence.Value) DisplayView {
	fields := value.(cadence.Struct).Fields
	return DisplayView{
		fields[0].ToGoValue().(string),
		fields[1].ToGoValue().(string),
		fields[2].ToGoValue().(string),
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
