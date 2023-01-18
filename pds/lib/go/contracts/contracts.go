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
)

const (
	filenameIPackNFT = "IPackHFT.cdc"
)

// IPackNFT returns the IPackNFT contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func IPackNFT(nftAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameIPackNFT)

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())

	return []byte(code)
}
