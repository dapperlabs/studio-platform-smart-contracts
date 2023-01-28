package templates

import (
	"github.com/onflow/flow-go-sdk"

	_ "github.com/kevinburke/go-bindata"

	"github.com/dapperlabs/studio-platform-smart-contracts/lib/go/templates/internal/assets"
)

const (
	filenameDeployPackNFT             = "transactions/deploy/deploy-packNFT-with-auth.cdc"
	filenameDeployPDS                 = "transactions/deploy/deploy-pds-with-auth.cdc"
	filenameCreatePackIssuer          = "transactions/pds/create_new_pack_issuer.cdc"
	filenameLinkExampleNFTProviderCap = "transactions/exampleNFT/link_providerCap_exampleNFT.cdc"
	filenmaeSetPackIssuerCap          = "transactions/pds/set_pack_issuer_cap.cdc"
	filenameCreateDistribution        = "transactions/pds/create_distribution.cdc"
	filenameMintPackNFT               = "transactions/pds/mint_packNFT.cdc"
)

// GenerateDeployPackNFTTx returns a transaction script that
// links a new royalty receiver interface
func GenerateDeployPackNFTTx(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameDeployPackNFT)
	return replaceAddresses(code, nftAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, iPackNFTAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateDeployPDSTx returns a transaction script that instantiates a new
// NFT collection instance, saves the collection in storage, then stores a
// reference to the collection.
func GenerateDeployPDSTx(nftAddress, iPackNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameDeployPDS)
	return replaceAddresses(code, nftAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, iPackNFTAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateCreatePackIssuerTx returns a transaction script that instantiates a new
// PackIssuer instance, saves it in storage, then stores a reference to it.
func GenerateCreatePackIssuerTx(pdsAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameCreatePackIssuer)
	return replaceAddresses(code, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, pdsAddress, flow.EmptyAddress)
}

// GenerateLinkExampleNFTProviderCapTx returns a transaction script that links NFT provider to a private path
func GenerateLinkExampleNFTProviderCapTx(nftAddress, exampleNFTAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameLinkExampleNFTProviderCap)
	return replaceAddresses(code, nftAddress, exampleNFTAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress)
}

// GenerateSetPackIssuerCapTx returns a transaction script that sets the pack issuer capability
func GenerateSetPackIssuerCapTx(pdsAddress flow.Address) []byte {
	code := assets.MustAssetString(filenmaeSetPackIssuerCap)
	return replaceAddresses(code, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, pdsAddress, flow.EmptyAddress)
}

// GenerateCreateDistributionTx returns a transaction script that creates a distribution
func GenerateCreateDistributionTx(pdsAddress, packNFTAddress, iPackNFTAddress, nftAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameCreateDistribution)
	return replaceAddresses(code, nftAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, iPackNFTAddress, pdsAddress, packNFTAddress)
}

// GenerateMintPackNFTTx returns a transaction script that mints a pack NFT
func GenerateMintPackNFTTx(pdsAddress, packNFTAddress, nftAddress flow.Address) []byte {
	code := assets.MustAssetString(filenameMintPackNFT)
	return replaceAddresses(code, nftAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, flow.EmptyAddress, pdsAddress, packNFTAddress)
}
