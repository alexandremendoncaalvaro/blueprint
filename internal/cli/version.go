package cli

import (
	"fmt"

	"github.com/ale/dotfiles/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Mostrar versao do dotfiles",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("dotfiles %s (commit: %s, data: %s)\n",
				version.Version, version.Commit, version.Date)
		},
	}
}
