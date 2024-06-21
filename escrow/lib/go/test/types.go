package test

import (
	"fmt"
	"log"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type SeriesData struct {
	ID     uint64
	Name   string
	Active bool
}
type SetData struct {
	ID   uint64
	Name string
}
type PlayData struct {
	ID             uint64
	Classification string
	Metadata       map[string]string
}
type EditionData struct {
	ID          uint64
	SeriesID    uint64
	SetID       uint64
	PlayID      uint64
	MaxMintSize *uint64
	Tier        string
}
type OurNFTData struct {
	ID           uint64
	EditionID    uint64
	SerialNumber uint64
	// A UFix64 in uint64 form
	MintingDate uint64
}
type LeaderboardInfo struct {
	Name          string
	NftType       string
	EntriesLength uint64
}

func cadenceStringDictToGo(cadenceDict cadence.Dictionary) map[string]string {
	goDict := make(map[string]string)
	for _, pair := range cadenceDict.Pairs {
		goDict[GetFieldValue(pair.Key).(string)] = GetFieldValue(pair.Value).(string)
	}
	return goDict
}

func parseSeriesData(value cadence.Value) SeriesData {
	s := SeriesData{}
	for k, v := range cadence.FieldsMappedByName(value.(cadence.Struct)) {
		switch k {
		case "id":
			s.ID = GetFieldValue(v).(uint64)
		case "name":
			s.Name = GetFieldValue(v).(string)
		case "active":
			s.Active = GetFieldValue(v).(bool)
		default:
			log.Fatalf("parseSeriesData: unexpected field: %s", k)
		}
	}
	return s
}

func parseSetData(value cadence.Value) SetData {
	s := SetData{}
	for k, v := range cadence.FieldsMappedByName(value.(cadence.Struct)) {
		switch k {
		case "id":
			s.ID = GetFieldValue(v).(uint64)
		case "name":
			s.Name = GetFieldValue(v).(string)
		case "setPlaysInEditions":
			continue
		default:
			log.Fatalf("parseSetData: unexpected field: %s", k)
		}
	}
	return s
}

func parsePlayData(value cadence.Value) PlayData {
	s := PlayData{}
	for k, v := range cadence.FieldsMappedByName(value.(cadence.Struct)) {
		switch k {
		case "id":
			s.ID = GetFieldValue(v).(uint64)
		case "classification":
			s.Classification = GetFieldValue(v).(string)
		case "metadata":
			s.Metadata = cadenceStringDictToGo(v.(cadence.Dictionary))
		default:
			log.Fatalf("parsePlayData: unexpected field: %s", k)
		}
	}
	return s
}

func parseEditionData(value cadence.Value) EditionData {
	s := EditionData{}
	for k, v := range cadence.FieldsMappedByName(value.(cadence.Struct)) {
		switch k {
		case "id":
			s.ID = GetFieldValue(v).(uint64)
		case "seriesID":
			s.SeriesID = GetFieldValue(v).(uint64)
		case "setID":
			s.SetID = GetFieldValue(v).(uint64)
		case "playID":
			s.PlayID = GetFieldValue(v).(uint64)
		case "maxMintSize":
			if f := GetFieldValue(v); f != nil {
				maxMintSize := f.(uint64)
				s.MaxMintSize = &maxMintSize
			}
		case "tier":
			s.Tier = GetFieldValue(v).(string)
		case "numMinted":
			continue
		default:
			log.Fatalf("parseEditionData: unexpected field: %s", k)
		}
	}
	return s
}

func parseNFTProperties(value cadence.Value) OurNFTData {
	array := value.(cadence.Array).Values
	return OurNFTData{
		GetFieldValue(array[0]).(uint64),
		GetFieldValue(array[1]).(uint64),
		GetFieldValue(array[2]).(uint64),
		GetFieldValue(array[3]).(uint64),
	}
}

func parseLeaderboardInfo(value cadence.Value) (LeaderboardInfo, error) {
	optionalVal, ok := value.(cadence.Optional)
	if !ok {
		return LeaderboardInfo{}, fmt.Errorf("expected value to be of type cadence.Optional, got %T", value)
	}

	if optionalVal.Value == nil {
		return LeaderboardInfo{}, fmt.Errorf("optional value is nil")
	}

	s := LeaderboardInfo{}
	fields := cadence.FieldsMappedByName(optionalVal.Value.(cadence.Struct))
	if len(fields) < 3 {
		return LeaderboardInfo{}, fmt.Errorf("struct does not contain enough fields")
	}
	for k, v := range fields {
		switch k {
		case "name":
			s.Name = v.(cadence.String).String()
		case "nftType":
			s.NftType = v.String()
		case "entriesLength":
			s.EntriesLength = v.(cadence.Int).Value.Uint64()
		default:
			log.Fatalf("parseLeaderboardInfo: unexpected field: %s", k)
		}
	}
	return s, nil
}

// GetFieldValue Convert a cadence value into a interface{} structure for easier consumption in go with options
func GetFieldValue(md cadence.Value) interface{} {
	switch field := md.(type) {
	case cadence.Optional:
		if field.Value == nil {
			return nil
		}
		return GetFieldValue(field.Value)
	case cadence.Dictionary:
		return convertDict(field)
	case cadence.Array:
		return convertArray(field)
	case cadence.Address:
		return flow.BytesToAddress(field.Bytes())
	case cadence.Int8:
		return int8(field)
	case cadence.Int16:
		return int16(field)
	case cadence.Int32:
		return int32(field)
	case cadence.Int64:
		return int64(field)
	case cadence.Int:
		return field.Value
	case cadence.UInt8:
		return uint8(field)
	case cadence.UInt16:
		return uint16(field)
	case cadence.UInt32:
		return uint32(field)
	case cadence.UInt64:
		return uint64(field)
	case cadence.Word8:
		return uint8(field)
	case cadence.Word16:
		return uint16(field)
	case cadence.Word32:
		return uint32(field)
	case cadence.Word64:
		return uint64(field)
	case cadence.TypeValue:
		return field.StaticType.ID()
	case cadence.String:
		return string(field)
	case cadence.UFix64:
		return uint64(field)
	case cadence.Fix64:
		return int64(field)
	case cadence.Struct:
		return ConvertObjectMetadata(field)
	case cadence.Resource:
		return ConvertObjectMetadata(field)
	case cadence.Bool:
		return bool(field)
	case cadence.Bytes:
		return []byte(field)
	case cadence.Character:
		return string(field)
	case cadence.Function:
		return field.FunctionType.ID()
	default:
		return field.String()
	}
}
func ConvertObjectMetadata(value cadence.Composite) map[string]interface{} {
	structMap := map[string]interface{}{}
	subFields := cadence.FieldsMappedByName(value)
	for key, subField := range subFields {
		if GetFieldValue(subField) != nil {
			structMap[key] = value
		}
	}
	return structMap
}

func convertArray(md cadence.Value) interface{} {
	arr := md.(cadence.Array)
	var items []interface{}
	for _, v := range arr.Values {
		items = append(items, GetFieldValue(v))
	}
	return items
}

func convertDict(md cadence.Value) interface{} {
	d := md.(cadence.Dictionary)
	valMap := map[string]interface{}{}
	for _, item := range d.Pairs {
		if item.Key.Type() != cadence.StringType {
			log.Fatalf("keys must be string type got %T", item.Key)
		}
		key := string(item.Key.(cadence.String))
		if key == "" {
			log.Fatalf("keys cannot be empty")
		}
		if v := GetFieldValue(item.Value); v != nil {
			valMap[key] = v
		}
	}
	return valMap
}
