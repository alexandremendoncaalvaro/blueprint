// Package moduletest fornece helpers de teste para o pacote module.
package moduletest

import "github.com/ale/dotfiles/internal/module"

// noopReporter implementa module.Reporter descartando todas as mensagens.
type noopReporter struct{}

func (r *noopReporter) Info(_ string)              {}
func (r *noopReporter) Success(_ string)           {}
func (r *noopReporter) Warn(_ string)              {}
func (r *noopReporter) Error(_ string)             {}
func (r *noopReporter) Step(_, _ int, _ string)    {}

// NoopReporter retorna um Reporter que descarta todas as mensagens.
// Util em testes onde o output do reporter nao importa.
func NoopReporter() module.Reporter {
	return &noopReporter{}
}
