package pds

import (
	"embed"
)

//go:embed transactions/*
var Transaction embed.FS

//go:embed scripts/*
var Scripts embed.FS

//go:embed contracts/*
var Contracts embed.FS
