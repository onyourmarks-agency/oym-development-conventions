package conventions

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/buildinfo"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/project"
)

type checkOptions struct {
	target string
	stacks string
}

func NewCheckCmd(assets fs.FS) *cobra.Command {
	o := &checkOptions{}
	c := &cobra.Command{
		Use:   "check",
		Short: "Exit non-zero when the vendored convention blocks are stale or missing (CI gate).",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runCheck(assets, o)
		},
	}
	c.Flags().StringVar(&o.target, "target", ".", "project directory")
	c.Flags().StringVar(&o.stacks, "stacks", "", "comma-separated stacks that must be present")
	return c
}

func runCheck(assets fs.FS, o *checkOptions) error {
	result, err := project.Sync(project.Options{
		Dir:       o.target,
		Assets:    assets,
		Stacks:    splitStacks(o.stacks),
		Rev:       buildinfo.Version,
		Date:      today(),
		CheckOnly: true,
	})
	if err != nil {
		return err
	}
	if len(result.Appended) > 0 {
		return fmt.Errorf("missing convention blocks: %v (run `oym-conventions sync --stacks ...`)", result.Appended)
	}
	if len(result.Updated) > 0 {
		return fmt.Errorf("stale convention blocks: %v (run `oym-conventions sync`)", result.Updated)
	}
	fmt.Println("cards up to date")
	printUpdateHint()
	return nil
}
