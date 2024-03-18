package templates

import (
	"regexp"

	"github.com/onflow/flow-go-sdk"
)

//go:generate go-bindata -prefix ../../../ -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../scripts/... ../../../transactions/...

var (
	placeholderNonFungibleToken  = regexp.MustCompile(`"NonFungibleToken"`)
	placeholderNonFungibleToken2 = regexp.MustCompile(`"NonFungibleToken"`)
	placeholderExampleNFT        = regexp.MustCompile(`"ExampleNFT"`)
	placeholderExampleNFT2       = regexp.MustCompile(`"ExampleNFT"`)
	placeHolderExampleNFT3       = regexp.MustCompile(`ExampleNFT`)
	placeHolderExampleNFT4       = regexp.MustCompile(`"ExampleNFT"`)
	placeholderMetadataViews     = regexp.MustCompile(`"MetadataViews"`)
	placeholderFungibleToken     = regexp.MustCompile(`"FungibleToken"`)
	placeholderIPackNFT          = regexp.MustCompile(`"IPackNFT"`)
	placeholderIPackNFT2         = regexp.MustCompile(`"IPackNFT"`)
	placeholderPackNFT           = regexp.MustCompile(`"PackNFT"`)
	placeholderPackNFT2          = regexp.MustCompile(`"PackNFT"`)
	placeholderPackNFT3          = regexp.MustCompile(`"PackNFT"`)
	placeholderPackNFT4          = regexp.MustCompile(`{{.PackNFTName}}`)
	placeholderPDS               = regexp.MustCompile(`"PDS"`)
)

func replaceAddresses(code string, nftAddress, exampleNFTAddress, metadataAddress, ftAddress, iPackNFTAddress, pdsAddress, packNFTAddress flow.Address) []byte {
	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderNonFungibleToken2.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderExampleNFT.ReplaceAllString(code, "0x"+exampleNFTAddress.String())
	code = placeholderExampleNFT2.ReplaceAllString(code, "0x"+exampleNFTAddress.String())
	code = placeHolderExampleNFT3.ReplaceAllString(code, "ExampleNFT")
	code = placeHolderExampleNFT4.ReplaceAllString(code, "0x"+exampleNFTAddress.String())
	code = placeholderMetadataViews.ReplaceAllString(code, "0x"+metadataAddress.String())
	code = placeholderFungibleToken.ReplaceAllString(code, "0x"+ftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())
	code = placeholderIPackNFT2.ReplaceAllString(code, "0x"+iPackNFTAddress.String())
	code = placeholderPackNFT.ReplaceAllString(code, "0x"+packNFTAddress.String())
	code = placeholderPackNFT2.ReplaceAllString(code, "0x"+packNFTAddress.String())
	code = placeholderPackNFT3.ReplaceAllString(code, "0x"+packNFTAddress.String())
	code = placeholderPackNFT4.ReplaceAllString(code, "PackNFT")
	code = placeholderPDS.ReplaceAllString(code, "0x"+pdsAddress.String())
	return []byte(code)
}
