package cli

import (
	"github.com/ale/blueprint/internal/module"
	"github.com/spf13/cobra"
)

// Options armazena as flags globais da CLI.
type Options struct {
	Headless bool
	Profile  string
	DryRun   bool
	Verbose  bool
}

// App agrupa as dependencias necessarias para os comandos.
type App struct {
	Registry  *module.Registry
	System    module.System
	Options   *Options
	ConfigDir string // Caminho para o diretorio configs/ do repo
}

// NewRootCmd cria o comando raiz com todas as flags globais.
func NewRootCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "blueprint",
		Short:         "Gerenciador de configuracoes para Bluefin",
		Long:          "CLI para configurar e manter seu ambiente Bluefin, com suporte a TUI interativo e modo headless.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Flags globais
	cmd.PersistentFlags().BoolVar(&app.Options.Headless, "headless", false, "Modo headless (sem TUI)")
	cmd.PersistentFlags().StringVarP(&app.Options.Profile, "profile", "p", "auto", "Perfil de instalacao (auto, full, minimal, server)")
	cmd.PersistentFlags().BoolVar(&app.Options.DryRun, "dry-run", false, "Mostrar o que seria feito sem executar")
	cmd.PersistentFlags().BoolVarP(&app.Options.Verbose, "verbose", "v", false, "Saida detalhada")

	// Subcomandos
	cmd.AddCommand(
		newApplyCmd(app),
		newStatusCmd(app),
		newUpdateCmd(app),
		newVersionCmd(),
	)

	return cmd
}
