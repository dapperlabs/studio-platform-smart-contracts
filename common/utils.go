package common

import "regexp"

var (
	placeholderNonFungibleToken = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	placeholderExampleNFT       = regexp.MustCompile(`"[^"\s].*/ExampleNFT.cdc"`)
	placeholderMetadataViews    = regexp.MustCompile(`"[^"\s].*/MetadataViews.cdc"`)
	placeholderFungibleToken    = regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)
)

type Replacer map[AddressName]*regexp.Regexp

type AddressName string

const (
	NonFungibleTokenAddr AddressName = "NonFungibleToken"
)

var NewReplacer Replacer

func init() {
	NewReplacer = map[AddressName]*regexp.Regexp{
		"placeholderNonFungibleToken": placeholderNonFungibleToken,
		"placeholderExampleNFT":       placeholderExampleNFT,
		"placeholderMetadataViews":    placeholderMetadataViews,
		"placeholderFungibleToken":    placeholderFungibleToken,
	}
}

func (r Replacer) replaceAddresses(code string, addr map[AddressName]string) string {
	for key, value := range addr {
		code = r[key].ReplaceAllString(code, "0x"+value)
	}
	return code
}
