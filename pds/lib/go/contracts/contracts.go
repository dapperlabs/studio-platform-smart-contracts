package contracts

//go:generate go-bindata -prefix ../../../contracts -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../contracts

import (
	//"regexp"

	"regexp"

	_ "github.com/kevinburke/go-bindata"

	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/contracts/internal/assets"
	"github.com/onflow/flow-go-sdk"
)

var (
	placeholderNonFungibleToken = regexp.MustCompile(`"NonFungibleToken"`)
	placeholderFungibleToken    = regexp.MustCompile(`"FungibleToken"`)
	placeholderIPackNFT         = regexp.MustCompile(`"IPackNFT"`)
	placeholderMetadataViews    = regexp.MustCompile(`"MetadataViews"`)
	placeholderRoyaltyAddress   = regexp.MustCompile(`{{.RoyaltyAddress}}`)
)

const (
	filenameIPackNFT      = "IPackNFT.cdc"
	filenamePackNFT       = "PackNFT.cdc"
	filenameAllDayPackNFT = "PackNFT_AllDay.cdc"
	filenamePDS           = "PDS.cdc"
)

// IPackNFT returns the IPackNFT contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func IPackNFT(nftAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameIPackNFT))

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	return []byte(code)
}

// PackNFT returns the PackNFT contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func PackNFT(nftAddress, iPackNFTAddress, metaDataViewAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenamePackNFT))
	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())
	code = placeholderMetadataViews.ReplaceAllString(code, "0x"+metaDataViewAddress.String())

	return []byte(code)
}

// AllDayPackNFT returns the AllDayPackNFT contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func AllDayPackNFT(nftAddress, ftAddress, iPackNFTAddress, metaDataViewAddress, packNFTAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameAllDayPackNFT))

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())
	code = placeholderFungibleToken.ReplaceAllString(code, "0x"+ftAddress.String())
	code = placeholderMetadataViews.ReplaceAllString(code, "0x"+metaDataViewAddress.String())
	code = placeholderRoyaltyAddress.ReplaceAllString(code, "0x"+packNFTAddress.String())

	return []byte(code)
}

// PDS returns the PDS contract.
//
// The returned contract will import the PDS contract from the specified address.
func PDS(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenamePDS))

	code = placeholderNonFungibleToken.ReplaceAllString(code, "0x"+nftAddress.String())
	code = placeholderIPackNFT.ReplaceAllString(code, "0x"+iPackNFTAddress.String())

	return []byte(code)
}
