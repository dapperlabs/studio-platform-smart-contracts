package templates

import (
	"github.com/onflow/flow-go-sdk"

	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates/internal/assets"
)

const (
	filenameBorrowNFT           = "scripts/exampleNFT/borrow_nft.cdc"
	filenameGetCollectionLength = "scripts/exampleNFT/get_collection_length.cdc"
	filenameGetTotalSupply      = "scripts/exampleNFT/get_total_supply.cdc"
	filenameGetNFTMetadata      = "scripts/get_nft_metadata.cdc"
	filenameGetNFTView          = "scripts/get_nft_view.cdc"
	filenameGetDistTitle        = "scripts/pds/get_dist_title.cdc"
	filenamePackNFTTotalSupply  = "scripts/packNFT/packNFT_total_supply.cdc"
)

// GenerateBorrowNFTScript creates a script that retrieves an NFT collection
// from storage and tries to borrow a reference for an NFT that it owns.
// If it owns it, it will not fail.
func GenerateBorrowNFTScript(nftAddress, exampleNFTAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameBorrowNFT))
	return replaceAddresses(code, nftAddress, exampleNFTAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateGetTotalSupplyScript creates a script that retrieves an NFT collection
// from storage and tries to borrow a reference for an NFT that it owns.
// If it owns it, it will not fail.
func GenerateGetTotalSupplyScript(nftAddress, exampleNFTAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameGetTotalSupply))
	return replaceAddresses(code, nftAddress, exampleNFTAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateGetNFTMetadataScript creates a script that returns the metadata for an NFT.
func GenerateGetNFTMetadataScript(nftAddress, exampleNFTAddress, metadataAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameGetNFTMetadata))
	return replaceAddresses(code, nftAddress, exampleNFTAddress, metadataAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateGetNFTViewScript creates a script that returns the rollup NFT View for an NFT.
func GenerateGetNFTViewScript(nftAddress, exampleNFTAddress, metadataAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameGetNFTView))
	return replaceAddresses(code, nftAddress, exampleNFTAddress, metadataAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateGetCollectionLengthScript creates a script that retrieves an NFT collection
// from storage and tries to borrow a reference for an NFT that it owns.
// If it owns it, it will not fail.
func GenerateGetCollectionLengthScript(nftAddress, exampleNFTAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameGetCollectionLength))
	return replaceAddresses(code, nftAddress, exampleNFTAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
}

//// GenerateGetTotalSupplyScript creates a script that reads
//// the total supply of tokens in existence
//// and makes assertions about the number
//func GenerateGetTotalSupplyScript(nftAddress, exampleNFTAddress flow.Address) []byte {
//	code := assets.MustAssetString(filenameGetTotalSupply)
//	return replaceAddresses(code, nftAddress, exampleNFTAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
//}

// GenerateGetDistTitleScript creates a script that returns the title of a distribution
func GenerateGetDistTitleScript(pdsAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenameGetDistTitle))
	return replaceAddresses(code, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, pdsAddress, flow.EmptyAddress)
}

// GenerateGetDistTitleScript creates a script that returns the total supply of pack NFTs
func GeneratePackNFTTotalSupply(exampleNFTAddress flow.Address) []byte {
	code := string(assets.MustAsset(filenamePackNFTTotalSupply))
	return replaceAddresses(code, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, exampleNFTAddress)
}
