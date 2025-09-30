package server

import "embed"

//go:embed "assets"
var Files embed.FS

//go:embed "locales"
var Locales embed.FS
