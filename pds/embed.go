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
	//go:embed transactions/packNFT/batch_transfer_packNFTs.cdc
	PackNFTBatchTransferPackNFT []byte
	//go:embed transactions/packNFT/transfer_packNFT.cdc
	PackNFTTransferPackNFT []byte
	//go:embed transactions/pds/mint_packNFT.cdc
	PDSMintPackNFT []byte
	//go:embed transactions/pds/open_packNFT.cdc
	PDSOpenPackNFT []byte
	//go:embed transactions/pds/reveal_packNFT.cdc
	PDSRevealPackNFT []byte
)
