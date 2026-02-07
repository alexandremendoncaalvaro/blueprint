package tui

import (
	"fmt"
	"strings"

	"github.com/ale/blueprint/internal/profile"
	tea "github.com/charmbracelet/bubbletea"
)

// profileSelectModel permite escolher o perfil de instalacao.
type profileSelectModel struct {
	profiles []profile.Profile
	cursor   int
	selected profile.Profile
	done     bool
}

func newProfileSelectModel() profileSelectModel {
	return profileSelectModel{
		profiles: profile.All(),
	}
}

func (m profileSelectModel) Init() tea.Cmd {
	return nil
}

func (m profileSelectModel) Update(msg tea.Msg) (profileSelectModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.profiles[m.cursor]
			m.done = true
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m profileSelectModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Selecione o perfil"))
	b.WriteString("\n\n")

	for i, p := range m.profiles {
		cursor := "  "
		style := mutedStyle
		if i == m.cursor {
			cursor = highlightStyle.Render("> ")
			style = highlightStyle
		}

		name := style.Render(p.Name)
		desc := mutedStyle.Render(fmt.Sprintf(" — %s", p.Description))
		tags := mutedStyle.Render(fmt.Sprintf(" [tags: %s]", strings.Join(p.Tags, ", ")))

		b.WriteString(cursor + name + desc + tags + "\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Use j/k ou setas para navegar, ENTER para selecionar"))

	return boxStyle.Render(b.String())
}
