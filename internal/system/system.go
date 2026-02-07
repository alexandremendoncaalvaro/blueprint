// Package system fornece implementacoes concretas de module.System.
// Inclui Real (SO), Mock (testes) e DryRun (simulacao).
package system

import "github.com/ale/dotfiles/internal/module"

// System e um alias para module.System.
// Mantido para compatibilidade; prefira usar module.System diretamente.
type System = module.System
