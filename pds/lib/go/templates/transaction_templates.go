package templates

import (
	"github.com/onflow/flow-go-sdk"

	_ "github.com/kevinburke/go-bindata"

	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates/internal/assets"
)

const (
	filenameDeployPackNFT = "transactions/deploy/deploy-packNFT-with-auth.cdc"
	filenameDeployPDS     = "transactions/deploy/deploy-pds-with-auth.cdc"
)

// GenerateDeployPackNFTTx returns a script that
// links a new royalty receiver interface
func GenerateDeployPackNFTTx(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameDeployPackNFT)
	return replaceAddresses(code, nftAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, iPackNFTAddress)
}

// GenerateDeployPDSTx returns a script that instantiates a new
// NFT collection instance, saves the collection in storage, then stores a
// reference to the collection.
func GenerateDeployPDSTx(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameDeployPDS)
	return replaceAddresses(code, nftAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, iPackNFTAddress)
}
