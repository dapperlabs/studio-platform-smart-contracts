package contracts

//go:generate go-bindata -prefix ../../../contracts -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../contracts

import (
	"regexp"

	_ "github.com/kevinburke/go-bindata"

	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/contracts/internal/assets"
	"github.com/onflow/flow-go-sdk"
)

var (
	placeholderNonFungibleToken = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	placeholderIPackNFT         = regexp.MustCompile(`"[^"\s].*/IPackNFT.cdc"`)
)

const (
	filenameIPackNFT = "IPackNFT.cdc"
	filenamePackNFT  = "PackNFT.cdc"
	filenamePDS      = "PDS.cdc"
)

// IPackNFT returns the IPackNFT contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func IPackNFT(nftAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameIPackNFT)

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())

	return []byte(code)
}

// PackNFT returns the PackNFT contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func PackNFT(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenamePackNFT)

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())

	return []byte(code)
}

// PDS returns the PDS contract.
//
// The returned contract will import the PDS contract from the specified address.
func PDS(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenamePDS)

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())

	return []byte(code)
}
