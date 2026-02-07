package cli

import (
	"fmt"

	"github.com/ale/dotfiles/internal/module"
	"github.com/ale/dotfiles/internal/orchestrator"
	"github.com/ale/dotfiles/internal/system"
	"github.com/ale/dotfiles/internal/tui"
	"github.com/spf13/cobra"
)

func newUpdateCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Atualizar sistema Bluefin (rpm-ostree, Flatpak, fwupd, Distrobox)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Busca o modulo bluefin-update diretamente
			mod, ok := app.Registry.ByName("bluefin-update")
			if !ok {
				return fmt.Errorf("modulo bluefin-update nao encontrado")
			}

			sys := app.System
			if app.Options.DryRun {
				sys = system.NewDryRun(app.System, func(msg string) {
					fmt.Println(msg)
				})
			}

			reporter := tui.NewHeadlessReporter()
			orch := orchestrator.New(sys, reporter)
			results := orch.Run(cmd.Context(), []module.Module{mod})

			if len(results) > 0 && results[0].Err != nil {
				return results[0].Err
			}

			return nil
		},
	}
}
