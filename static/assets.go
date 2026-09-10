// Package static embeds the interface so releases run independently of the working directory.
package static

import "embed"

//go:embed pages/*.html styles/* js/* images/*
var Files embed.FS
