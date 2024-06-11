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
