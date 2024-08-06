package test

import (
	"fmt"
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
	fields := value.(cadence.Struct).FieldsMappedByName()
	metadata, err := convertDict(fields["metadata"])
	if err != nil {
		panic(fmt.Errorf("error converting metadata: %w", err))
	}
	md := make(map[string]string)
	for k, v := range metadata {
		md[k.(string)] = v.(string)
	}
	return EditionData{
		uint64(fields["id"].(cadence.UInt64)),
		uint64(fields["numMinted"].(cadence.UInt64)),
		bool(fields["active"].(cadence.Bool)),
		md,
	}
}

func parseEditionNftProperties(value cadence.Value) NFTData {
	array := value.(cadence.Array).Values
	return NFTData{
		uint64(array[0].(cadence.UInt64)),
		uint64(array[1].(cadence.UInt64)),
	}
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
