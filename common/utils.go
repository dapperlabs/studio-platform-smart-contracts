package common

import (
	"fmt"
	"reflect"
	"regexp"
)

var (
	placeholderNonFungibleToken = regexp.MustCompile(`"[^"\s]?.*NonFungibleToken(\.cdc)?"`)
	placeholderExampleNFT       = regexp.MustCompile(`"[^"\s].*ExampleNFT.cdc"`)
	placeholderMetadataViews    = regexp.MustCompile(`"[^"\s].*MetadataViews.cdc"`)
	placeholderFungibleToken    = regexp.MustCompile(`"[^"\s].*FungibleToken.cdc"`)
)

type Replacer map[AddressName]*regexp.Regexp

type AddressName string

const (
	NonFungibleTokenAddr AddressName = "NonFungibleToken"
)

func NewReplacer() Replacer {
	xxx := map[AddressName]*regexp.Regexp{
		NonFungibleTokenAddr: placeholderNonFungibleToken,
	}
	return xxx
}

func (r Replacer) ReplaceAddresses(code string, data interface{}) (string, error) {

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("data must be a struct")
	}

	var code1 = code
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Name
		fieldValue := val.Field(i).String()

		placeholder := AddressName(fieldName)
		if r[placeholder] == nil {
			return "", fmt.Errorf("no placeholder found for %s", placeholder)
		}

		code1 = r[placeholder].ReplaceAllString(code1, "0x"+fieldValue)
		fmt.Printf(`"%s  %s"\n`, fieldName, fieldValue)
	}
	return code1, nil
}
