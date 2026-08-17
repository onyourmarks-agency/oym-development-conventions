package conventions

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/buildinfo"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/cards"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/project"
)

type syncOptions struct {
	target string
	stacks string
	mode   string
}

func NewSyncCmd(assets fs.FS) *cobra.Command {
	o := &syncOptions{}
	c := &cobra.Command{
		Use:   "sync",
		Short: "Refresh the vendored convention blocks in AGENTS.md.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSync(assets, o)
		},
	}
	c.Flags().StringVar(&o.target, "target", ".", "project directory")
	c.Flags().StringVar(&o.stacks, "stacks", "", "comma-separated stacks to append when missing")
	c.Flags().StringVar(&o.mode, "mode", cards.DefaultMode, "mode for newly appended blocks")
	return c
}

func runSync(assets fs.FS, o *syncOptions) error {
	result, err := project.Sync(project.Options{
		Dir:    o.target,
		Assets: assets,
		Stacks: splitStacks(o.stacks),
		Mode:   o.mode,
		Rev:    buildinfo.Version,
		Date:   today(),
	})
	if err != nil {
		return err
	}
	switch {
	case result.UpToDate():
		fmt.Println("cards up to date, nothing written")
	default:
		fmt.Printf("synced %s", result.AgentsPath)
		if len(result.Appended) > 0 {
			fmt.Printf(" (added %v)", result.Appended)
		}
		if len(result.Updated) > 0 {
			fmt.Printf(" (updated %v)", result.Updated)
		}
		fmt.Println()
	}
	printUpdateHint()
	return nil
}
