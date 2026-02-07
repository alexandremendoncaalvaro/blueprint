package tui

import (
	"fmt"
	"strings"

	"github.com/ale/blueprint/internal/module"
	"github.com/ale/blueprint/internal/profile"
	tea "github.com/charmbracelet/bubbletea"
)

// moduleEntry agrupa um modulo e seu estado de selecao.
type moduleEntry struct {
	mod      module.Module
	selected bool
}

// moduleConfirmModel permite ativar/desativar modulos individualmente.
type moduleConfirmModel struct {
	entries      []moduleEntry
	cursor       int
	done         bool
	profile      profile.Profile
	autoDetected bool
}

func newModuleConfirmModel(modules []module.Module) moduleConfirmModel {
	entries := make([]moduleEntry, len(modules))
	for i, m := range modules {
		entries[i] = moduleEntry{mod: m, selected: true}
	}
	return moduleConfirmModel{entries: entries}
}

func (m moduleConfirmModel) Init() tea.Cmd {
	return nil
}

func (m moduleConfirmModel) Update(msg tea.Msg) (moduleConfirmModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case " ", "x":
			m.entries[m.cursor].selected = !m.entries[m.cursor].selected
		case "enter":
			m.done = true
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m moduleConfirmModel) View() string {
	var b strings.Builder

	if m.autoDetected {
		b.WriteString(titleStyle.Render("Blueprint"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("Perfil detectado: %s", highlightStyle.Render(m.profile.Name)))
		b.WriteString(mutedStyle.Render(fmt.Sprintf(" — %s", m.profile.Description)))
		b.WriteString("\n\n")
		b.WriteString("Modulos selecionados:\n\n")
	} else {
		b.WriteString(titleStyle.Render("Confirmar modulos"))
		b.WriteString("\n\n")
	}

	for i, e := range m.entries {
		cursor := "  "
		if i == m.cursor {
			cursor = highlightStyle.Render("> ")
		}

		check := "[ ]"
		if e.selected {
			check = successStyle.Render("[x]")
		}

		name := e.mod.Name()
		desc := mutedStyle.Render(fmt.Sprintf(" — %s", e.mod.Description()))
		tags := mutedStyle.Render(fmt.Sprintf(" [%s]", strings.Join(e.mod.Tags(), ", ")))

		b.WriteString(fmt.Sprintf("%s%s %s%s%s\n", cursor, check, name, desc, tags))
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("ESPACO para toggle, ENTER para confirmar, q para sair"))

	return boxStyle.Render(b.String())
}

func (m moduleConfirmModel) selectedModules() []module.Module {
	var result []module.Module
	for _, e := range m.entries {
		if e.selected {
			result = append(result, e.mod)
		}
	}
	return result
}
