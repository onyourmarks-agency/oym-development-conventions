package conventions

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/buildinfo"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information.",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(buildinfo.String())
		},
	}
}
