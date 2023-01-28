package templates

import "regexp"

//go:generate go-bindata -prefix ../../../ -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../scripts/... ../../../transactions/...

import (
	"github.com/onflow/flow-go-sdk"
)

var (
	placeholderNonFungibleToken  = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	placeholderNonFungibleToken2 = regexp.MustCompile(`0x{{.NonFungibleToken}}`)
	placeholderExampleNFT        = regexp.MustCompile(`"[^"\s].*/ExampleNFT.cdc"`)
	placeholderExampleNFT2       = regexp.MustCompile(`0x{{.ExampleNFT}}`)
	placeHolderExampleNFT3		= regexp.MustCompile(`{{.CollectibleNFTName}}`)
	placeHolderExampleNFT4		= regexp.MustCompile(`0x{{.CollectibleNFTAddress}}`)
	placeholderMetadataViews     = regexp.MustCompile(`"[^"\s].*/MetadataViews.cdc"`)
	placeholderFungibleToken     = regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)
	placeholderIPackNFT          = regexp.MustCompile(`"[^"\s].*/IPackNFT.cdc"`)
	placeholderIPackNFT2          = regexp.MustCompile(`0x{{.IPackNFT}}`)
	placeholderPackNFT           = regexp.MustCompile(`"[^"\s].*/PackNFT.cdc"`)
	placeholderPackNFT2          = regexp.MustCompile(`0x{{.PackNFTAddress}}`)
	placeholderPackNFT3          = regexp.MustCompile(`0x{{.PackNFT}}`)
	placeholderPackNFT4          = regexp.MustCompile(`{{.PackNFTName}}`)
	placeholderPDS               = regexp.MustCompile(`0x{{.PDS}}`)

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
