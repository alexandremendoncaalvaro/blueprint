// Package version armazena informacoes de versao injetadas via ldflags.
package version

// Variaveis populadas em build time via ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
