package test

import (
	"log"
	"math/big"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type SeriesData struct {
	ID         *big.Int
	Name       string
	LockedDate *uint64
}

func cadenceAddressToBigIntDictToGo(cadenceDict cadence.Dictionary) map[flow.Address]*big.Int {
	goDict := make(map[flow.Address]*big.Int)
	for _, pair := range cadenceDict.Pairs {
		goDict[GetFieldValue(pair.Key).([8]uint8)] = GetFieldValue(pair.Value).(*big.Int)
	}
	return goDict
}

func cadenceStringArrayToGo(cadenceArray cadence.Array) []string {
	goArray := []string{}
	for _, i := range cadenceArray.Values {
		goArray = append(goArray, GetFieldValue(i).(string))
	}
	return goArray
}

func parseSeriesData(value cadence.Value) SeriesData {
	s := SeriesData{}
	for k, v := range cadence.FieldsMappedByName(value.(cadence.Struct)) {
		switch k {
		case "id":
			s.ID = GetFieldValue(v).(*big.Int)
		case "name":
			s.Name = GetFieldValue(v).(string)
		case "lockedDate":
			if f := GetFieldValue(v); f != nil {
				lockedData := f.(uint64)
				s.LockedDate = &lockedData
			}
		default:
			log.Fatalf("unexpected field: %s", k)
		}
	}
	return s
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
