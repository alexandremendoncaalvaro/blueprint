package cli

import (
	"fmt"

	"github.com/ale/blueprint/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Mostrar versao do blueprint",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("blueprint %s (commit: %s, data: %s)\n",
				version.Version, version.Commit, version.Date)
		},
	}
}
