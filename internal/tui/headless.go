package tui

import "fmt"

// HeadlessReporter implementa Reporter para modo headless (sem TUI).
// Usa saida textual simples com prefixos.
type HeadlessReporter struct{}

// NewHeadlessReporter cria um reporter para modo headless.
func NewHeadlessReporter() *HeadlessReporter {
	return &HeadlessReporter{}
}

func (r *HeadlessReporter) Info(msg string) {
	fmt.Printf("  [INFO] %s\n", msg)
}

func (r *HeadlessReporter) Success(msg string) {
	fmt.Printf("  [OK]   %s\n", msg)
}

func (r *HeadlessReporter) Warn(msg string) {
	fmt.Printf("  [WARN] %s\n", msg)
}

func (r *HeadlessReporter) Error(msg string) {
	fmt.Printf("  [ERRO] %s\n", msg)
}

func (r *HeadlessReporter) Step(current, total int, msg string) {
	fmt.Printf("  [%d/%d] %s\n", current, total, msg)
}
