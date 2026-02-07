package tui

import (
	"context"
	"strings"

	"github.com/ale/dotfiles/internal/module"
	"github.com/ale/dotfiles/internal/orchestrator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// executeModel mostra o progresso da execucao com spinner.
type executeModel struct {
	modules []module.Module
	sys     module.System
	spinner spinner.Model
	results []orchestrator.Result
	logs    []string
	done    bool
}

func newExecuteModel(modules []module.Module, sys module.System) executeModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = highlightStyle

	return executeModel{
		modules: modules,
		sys:     sys,
		spinner: s,
	}
}

func (m executeModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.runModules(),
	)
}

func (m executeModel) Update(msg tea.Msg) (executeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case modulesDoneMsg:
		m.results = msg.results
		m.done = true
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m executeModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Executando..."))
	b.WriteString("\n\n")

	if !m.done {
		b.WriteString(m.spinner.View())
		b.WriteString(" Aplicando configuracoes...\n\n")

		for _, mod := range m.modules {
			b.WriteString(mutedStyle.Render("  > " + mod.Name() + "\n"))
		}
	}

	return boxStyle.Render(b.String())
}

func (m executeModel) runModules() tea.Cmd {
	return func() tea.Msg {
		reporter := &tuiReporter{}
		orch := orchestrator.New(m.sys, reporter)
		results := orch.Run(context.Background(), m.modules)
		return modulesDoneMsg{results: results}
	}
}
