// Package tui implementa a interface grafica com Bubble Tea.
package tui

import "github.com/charmbracelet/lipgloss"

// Cores do tema.
var (
	colorPrimary   = lipgloss.Color("#7C3AED") // Roxo
	colorSuccess   = lipgloss.Color("#10B981") // Verde
	colorWarning   = lipgloss.Color("#F59E0B") // Amarelo
	colorError     = lipgloss.Color("#EF4444") // Vermelho
	colorMuted     = lipgloss.Color("#6B7280") // Cinza
	colorHighlight = lipgloss.Color("#3B82F6") // Azul
)

// Estilos reutilizaveis.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess)

	warningStyle = lipgloss.NewStyle().
			Foreground(colorWarning)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	highlightStyle = lipgloss.NewStyle().
			Foreground(colorHighlight).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2)
)
