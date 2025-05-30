package pds

import (
	"embed"
)

//go:embed transactions/*
var Transactions embed.FS

//go:embed scripts/*
var Scripts embed.FS

//go:embed contracts/*
var Contracts embed.FS

var (
	// ----- Scripts ----------------------------

	//go:embed scripts/packNFT/has_packNFT_collection.cdc
	HasPackNFTCollection []byte

	// ----- Transactions -----------------------

	//go:embed transactions/packNFT/batch_transfer_packNFTs.cdc
	PackNFTBatchTransferPackNFT []byte
	//go:embed transactions/packNFT/transfer_packNFT.cdc
	PackNFTTransferPackNFT []byte
	//go:embed transactions/pds/mint_packNFT.cdc
	PDSMintPackNFT []byte
	//go:embed transactions/pds/mint_packNFTs.cdc
	PDSMintPackNFTs []byte
	//go:embed transactions/pds/open_packNFT.cdc
	PDSOpenPackNFT []byte
	//go:embed transactions/pds/reveal_packNFT.cdc
	PDSRevealPackNFT []byte
	//go:embed transactions/pds/create_distribution.cdc
	PDSCreateDistribution []byte
	//go:embed transactions/pds/update_distribution.cdc
	PDSUpdateDistribution []byte
)
