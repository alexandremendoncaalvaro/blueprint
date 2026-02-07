// Package cli implementa os comandos da interface de linha de comando.
package cli

import (
	"os"

	"golang.org/x/term"
)

// Mode define o modo de operacao da CLI.
type Mode int

const (
	Interactive Mode = iota // TUI com Bubble Tea
	Headless               // Saida textual para automacao
)

// DetectMode decide automaticamente o modo de operacao.
// Usa headless se nao for terminal, se estiver em container, ou se --headless.
func DetectMode(forceHeadless bool) Mode {
	if forceHeadless {
		return Headless
	}

	// Se stdin nao e terminal, usa headless
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return Headless
	}

	// Se TERM nao esta definido, provavelmente nao e terminal
	if os.Getenv("TERM") == "" {
		return Headless
	}

	// CI/CD
	if os.Getenv("CI") != "" {
		return Headless
	}

	return Interactive
}
