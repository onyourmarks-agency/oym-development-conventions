package conventions

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/buildinfo"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/selfupdate"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Self-update to the latest GitHub release.",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			return selfupdate.Run(ctx, buildinfo.Version)
		},
	}
}
