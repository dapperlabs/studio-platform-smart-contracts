package templates

import "regexp"

//go:generate go-bindata -prefix ../../../ -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../scripts/... ../../../transactions/...

import (
	"github.com/onflow/flow-go-sdk"
)

var (
	placeholderNonFungibleToken = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	placeholderExampleNFT       = regexp.MustCompile(`"[^"\s].*/ExampleNFT.cdc"`)
	placeholderMetadataViews    = regexp.MustCompile(`"[^"\s].*/MetadataViews.cdc"`)
	placeholderFungibleToken    = regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)
	placeholderIPackNFT         = regexp.MustCompile(`"[^"\s].*/IPackNFT.cdc"`)
)

func replaceAddresses(code string, nftAddress, exampleNFTAddress, metadataAddress, ftAddress, iPackNFTAddress flow.Address) []byte {
	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderExampleNFT.ReplaceAllString(code, "0x"+exampleNFTAddress.String())
	code = placeholderMetadataViews.ReplaceAllString(code, "0x"+metadataAddress.String())
	code = placeholderFungibleToken.ReplaceAllString(code, "0x"+ftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())
	return []byte(code)
}
