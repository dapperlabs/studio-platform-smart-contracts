package test

import (
	"github.com/onflow/cadence"
)

type Series struct {
	ID     uint64
	Name   string
	Active bool
}

type Set struct {
	ID     uint64
	Name   string
	Locked bool
}

type Tag struct {
	ID   uint64
	Name string
}

type Play struct {
	ID       uint64
	Metadata map[string]string
	TagIds   []uint64
}

type Edition struct {
	ID          uint64
	SeriesID    uint64
	SetID       uint64
	PlayID      uint64
	MaxMintSize *uint64
	Tier        string
	NumMinted   uint64
}

type NFTData struct {
	ID           uint64
	EditionID    uint64
	SerialNumber uint64
}

type DisplayView struct {
	Name        string
	Description string
	ImageURL    string
}

func parseSeries(value cadence.Value) Series {
	fields := value.(cadence.Struct).Fields
	return Series{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(string),
		fields[2].ToGoValue().(bool),
	}
}

func parseSet(value cadence.Value) Set {
	fields := value.(cadence.Struct).Fields
	return Set{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(string),
		fields[2].ToGoValue().(bool),
	}
}

func parseTag(value cadence.Value) Tag {
	fields := value.(cadence.Struct).Fields
	return Tag{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(string),
	}
}

func parsePlay(value cadence.Value) Play {
	fields := value.(cadence.Struct).Fields
	var tagIds []uint64
	for _, val := range fields[2].(cadence.Array).Values {
		tagId := val.ToGoValue().(uint64)
		tagIds = append(tagIds, tagId)
	}
	return Play{
		fields[0].ToGoValue().(uint64),
		cadenceStringDictToGo(fields[1].(cadence.Dictionary)),
		tagIds,
	}
}

func parseEdition(value cadence.Value) Edition {
	fields := value.(cadence.Struct).Fields
	var maxMintSize uint64
	if fields[4] != nil && fields[4].ToGoValue() != nil {
		maxMintSize = fields[4].ToGoValue().(uint64)
	}
	return Edition{
		fields[0].ToGoValue().(uint64),
		fields[1].ToGoValue().(uint64),
		fields[2].ToGoValue().(uint64),
		fields[3].ToGoValue().(uint64),
		&maxMintSize,
		fields[5].ToGoValue().(string),
		fields[6].ToGoValue().(uint64),
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

func parseNFTProperties(value cadence.Value) NFTData {
	array := value.(cadence.Array).Values
	return NFTData{
		array[0].ToGoValue().(uint64),
		array[1].ToGoValue().(uint64),
		array[2].ToGoValue().(uint64),
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
