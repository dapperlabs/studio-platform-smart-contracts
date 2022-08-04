package test

import (
	"github.com/onflow/cadence"
)

type EditionData struct {
	ID        uint64
	NumMinted uint64
	Active    bool
	Metadata  map[string]string
}

type NFTData struct {
	ID        uint64
	EditionID uint64
}

func parseEditionData(value cadence.Value) EditionData {
	fields := value.(cadence.Struct).Fields
	return EditionData{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(uint64),
		fields[2].ToGoValue().(bool),
		cadenceStringDictToGo(fields[3].(cadence.Dictionary)),
	}
}

func parseEditionNftProperties(value cadence.Value) NFTData {
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
