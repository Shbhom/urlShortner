package server

import "embed"

//go:embed all:ui/dist
var embeddedUI embed.FS
