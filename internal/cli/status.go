package cli

import (
	"fmt"

	"github.com/ale/dotfiles/internal/module"
	"github.com/ale/dotfiles/internal/orchestrator"
	"github.com/ale/dotfiles/internal/profile"
	"github.com/ale/dotfiles/internal/tui"
	"github.com/spf13/cobra"
)

func newStatusCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Mostrar estado dos modulos",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var prof profile.Profile
			if app.Options.Profile == "auto" {
				prof = profile.Detect(app.System)
			} else {
				var err error
				prof, err = profile.ByName(app.Options.Profile)
				if err != nil {
					return err
				}
			}
			modules := profile.Resolve(prof, app.Registry)

			reporter := tui.NewHeadlessReporter()
			orch := orchestrator.New(app.System, reporter)
			results := orch.CheckAll(cmd.Context(), modules)

			fmt.Printf("Perfil: %s\n", prof.Name)
			fmt.Println()

			for _, r := range results {
				icon := statusIcon(r.Status.Kind)
				fmt.Printf("  %s %s — %s\n", icon, r.Module.Name(), r.Status.Message)
			}

			return nil
		},
	}
}

func statusIcon(kind module.StatusKind) string {
	switch kind {
	case module.Installed:
		return "[OK]   "
	case module.Missing:
		return "[FALTA]"
	case module.Partial:
		return "[PARC] "
	case module.Skipped:
		return "[SKIP] "
	default:
		return "[?]    "
	}
}
