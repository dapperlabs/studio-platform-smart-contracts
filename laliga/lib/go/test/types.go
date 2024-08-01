package test

import (
	"fmt"
	"github.com/onflow/cadence"
)

type SeriesData struct {
	ID     uint64
	Name   string
	Active bool
}
type SetData struct {
	ID     uint64
	Name   string
	Locked bool
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

type DisplayView struct {
	Name        string
	Description string
	ImageURL    string
}

type EditionView struct {
	Number uint64
	Max    uint64
}

type NFTCollectionDataView struct {
	StoragePath        string
	PublicPath         string
	ProviderPath       string
	PublicCollection   string
	PublicLinkedType   string
	ProviderLinkedType string
}

type TraitView struct {
	Name  string
	Value string
}

type TraitsView []TraitView

func cadenceStringDictToGo(cadenceDict cadence.Dictionary) map[string]string {
	goDict := make(map[string]string)
	for _, pair := range cadenceDict.Pairs {
		goDict[string(pair.Key.(cadence.String))] = string(pair.Value.(cadence.String))
	}
	return goDict
}

func parseSeriesData(value cadence.Value) SeriesData {
	fields := value.(cadence.Struct).FieldsMappedByName()
	return SeriesData{
		uint64(fields["id"].(cadence.UInt64)),
		string(fields["name"].(cadence.String)),
		bool(fields["active"].(cadence.Bool)),
	}
}

func parseSetData(value cadence.Value) SetData {
	fields := value.(cadence.Struct).FieldsMappedByName()
	return SetData{
		uint64(fields["id"].(cadence.UInt64)),
		string(fields["name"].(cadence.String)),
		bool(fields["locked"].(cadence.Bool)),
	}
}

func parsePlayData(value cadence.Value) PlayData {
	fields := value.(cadence.Struct).FieldsMappedByName()
	return PlayData{
		uint64(fields["id"].(cadence.UInt64)),
		string(fields["classification"].(cadence.String)),
		cadenceStringDictToGo(fields["metadata"].(cadence.Dictionary)),
	}
}

func parseEditionData(value cadence.Value) EditionData {
	fields := value.(cadence.Struct).FieldsMappedByName()
	var maxMintSize uint64
	if fields["maxMintSize"].(cadence.Optional).Value != nil {
		maxMintSize = uint64(fields["maxMintSize"].(cadence.Optional).Value.(cadence.UInt64))
	}
	return EditionData{
		uint64(fields["id"].(cadence.UInt64)),
		uint64(fields["seriesID"].(cadence.UInt64)),
		uint64(fields["setID"].(cadence.UInt64)),
		uint64(fields["playID"].(cadence.UInt64)),
		&maxMintSize,
		string(fields["tier"].(cadence.String)),
	}
}

func parseNFTProperties(value cadence.Value) OurNFTData {
	array := value.(cadence.Array).Values
	return OurNFTData{
		uint64(array[0].(cadence.UInt64)),
		uint64(array[1].(cadence.UInt64)),
		uint64(array[2].(cadence.UInt64)),
		uint64(array[3].(cadence.UFix64)),
	}
}

func parseMetadataDisplayView(value cadence.Value) DisplayView {
	fields := value.(cadence.Struct).FieldsMappedByName()
	return DisplayView{
		string(fields["name"].(cadence.String)),
		string(fields["description"].(cadence.String)),
		string(fields["thumbnail"].(cadence.String)),
	}
}

func parseMetadataEditionView(value cadence.Value) EditionView {
	fields := value.(cadence.Struct).FieldsMappedByName()
	edition := fields["infoList"].(cadence.Array).Values[0].(cadence.Struct).FieldsMappedByName()
	maxMintSize := uint64(0)
	if edition["max"].(cadence.Optional).Value != nil {
		maxMintSize = uint64(edition["max"].(cadence.Optional).Value.(cadence.UInt64))
	}
	return EditionView{
		uint64(edition["number"].(cadence.UInt64)),
		maxMintSize,
	}
}

func parseMetadataSerialView(value cadence.Value) uint64 {
	return uint64(value.(cadence.UInt64))
}

func parseMetadataNFTCollectionDataView(value cadence.Value) NFTCollectionDataView {
	fields := value.(cadence.Struct).FieldsMappedByName()
	return NFTCollectionDataView{
		StoragePath:      fields["storagePath"].(cadence.Path).Identifier,
		PublicPath:       fields["publicPath"].(cadence.Path).Identifier,
		PublicCollection: fields["publicCollection"].(cadence.TypeValue).StaticType.ID(),
		PublicLinkedType: fields["publicLinkedType"].(cadence.TypeValue).StaticType.ID(),
	}
}

func parseMetadataTraitsView(value cadence.Value) TraitsView {
	var view TraitsView
	fields := value.(cadence.Struct).FieldsMappedByName()
	for _, val := range fields["traits"].(cadence.Array).Values {
		trait := val.(cadence.Struct).FieldsMappedByName()
		traitVal, err := GetFieldValue(trait["value"])
		if err != nil {
			panic(err)
		}
		view = append(view, TraitView{
			Name:  string(trait["name"].(cadence.String)),
			Value: fmt.Sprintf("%v", traitVal),
		})
	}
	return view
}

func convertArray(md cadence.Value) (any, error) {
	arr := md.(cadence.Array)
	var items []any
	for _, v := range arr.Values {
		converted, err := GetFieldValue(v)
		if err != nil {
			return nil, err
		}
		items = append(items, converted)
	}
	return items, nil
}

func convertDict(md cadence.Value) (map[any]any, error) {
	d, ok := md.(cadence.Dictionary)
	if !ok {
		return nil, fmt.Errorf("value is not a dictionary, got %T", md)
	}
	valMap := map[any]any{}
	for _, item := range d.Pairs {
		value, err := GetFieldValue(item.Value)
		if err != nil {
			return nil, err
		}
		key, err := GetFieldValue(item.Key)
		if err != nil {
			return nil, err
		}
		if key == "" {
			return nil, fmt.Errorf("keys cannot be empty")
		}
		if value != nil {
			valMap[key] = value
		}
	}
	return valMap, nil
}

func ConvertObjectMetadata(value cadence.Composite) (map[string]any, error) {
	structMap := map[string]any{}
	subFields := cadence.FieldsMappedByName(value)
	for key, subField := range subFields {
		val, err := GetFieldValue(subField)
		if err != nil {
			return nil, err
		}
		if val != nil {
			structMap[key] = val
		}
	}
	return structMap, nil
}

func GetFieldValue(md cadence.Value) (any, error) {
	switch field := md.(type) {
	case cadence.Optional:
		if field.Value == nil {
			return nil, nil
		}
		return GetFieldValue(field.Value)
	case cadence.Dictionary:
		return convertDict(field)
	case cadence.Array:
		return convertArray(field)
	case cadence.Int:
		return field.Int(), nil
	case cadence.Int8:
		return int8(field), nil
	case cadence.Int16:
		return int16(field), nil
	case cadence.Int32:
		return int32(field), nil
	case cadence.Int64:
		return int64(field), nil
	case cadence.UInt8:
		return uint8(field), nil
	case cadence.UInt16:
		return uint16(field), nil
	case cadence.UInt32:
		return uint32(field), nil
	case cadence.UInt64:
		return uint64(field), nil
	case cadence.Word8:
		return uint8(field), nil
	case cadence.Word16:
		return uint16(field), nil
	case cadence.Word32:
		return uint32(field), nil
	case cadence.Word64:
		return uint64(field), nil
	case cadence.TypeValue:
		return field.StaticType.ID(), nil
	case cadence.String:
		return string(field), nil
	case cadence.UFix64:
		return uint64(field), nil
	case cadence.Fix64:
		return int64(field), nil
	case cadence.Struct:
		return ConvertObjectMetadata(field)
	case cadence.Resource:
		return ConvertObjectMetadata(field)
	case cadence.Bool:
		return bool(field), nil
	case cadence.Bytes:
		return []byte(field), nil
	case cadence.Character:
		return string(field), nil
	case cadence.Function:
		return field.FunctionType.ID(), nil
	case cadence.Address:
		return field.String()[2:], nil
	default:
		return field.String(), nil
	}
}
