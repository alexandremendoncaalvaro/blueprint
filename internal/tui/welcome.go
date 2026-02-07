package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// welcomeModel e a tela de boas-vindas.
type welcomeModel struct {
	done bool
}

func newWelcomeModel() welcomeModel {
	return welcomeModel{}
}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (welcomeModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", " ":
			m.done = true
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m welcomeModel) View() string {
	var b strings.Builder

	title := titleStyle.Render("Blueprint")
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString(subtitleStyle.Render("Configuracao automatizada do seu ambiente"))
	b.WriteString("\n\n")

	b.WriteString("Este assistente vai configurar:\n\n")
	b.WriteString(highlightStyle.Render("  > "))
	b.WriteString("Starship (prompt do shell)\n")
	b.WriteString(highlightStyle.Render("  > "))
	b.WriteString("Cedilha (fix para teclado US International)\n")
	b.WriteString(highlightStyle.Render("  > "))
	b.WriteString("Atualizacoes do sistema (rpm-ostree, Flatpak, etc.)\n")
	b.WriteString("\n\n")

	b.WriteString(mutedStyle.Render("Pressione ENTER para continuar ou q para sair"))

	return boxStyle.Render(b.String())
}
