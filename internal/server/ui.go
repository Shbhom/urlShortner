package server

import "embed"

//go:embed ui/dist/*
var embeddedUI embed.FS
