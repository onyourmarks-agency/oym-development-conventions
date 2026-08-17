// Package conventions wires the oym-conventions CLI.
package conventions

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/buildinfo"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/selfupdate"
)

func New(assets fs.FS) *cobra.Command {
	c := &cobra.Command{
		Use:   "oym-conventions",
		Short: "Set up and sync On Your Marks development conventions in a project.",
	}
	c.SilenceUsage = true
	c.SilenceErrors = true
	c.CompletionOptions.DisableDefaultCmd = true
	cobra.EnableCommandSorting = false

	c.AddCommand(NewInitCmd(assets))
	c.AddCommand(NewSyncCmd(assets))
	c.AddCommand(NewCheckCmd(assets))
	c.AddCommand(NewUpdateCmd())
	c.AddCommand(NewVersionCmd())
	return c
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func splitStacks(flagValue string) []string {
	if flagValue == "" {
		return nil
	}
	stacks := []string{}
	for _, stack := range strings.Split(flagValue, ",") {
		stack = strings.TrimSpace(stack)
		if stack != "" {
			stacks = append(stacks, stack)
		}
	}
	return stacks
}

func stdinIsTerminal() bool {
	fd := os.Stdin.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func printUpdateHint() {
	if tag := selfupdate.NewerAvailable(buildinfo.Version); tag != "" {
		fmt.Printf("note: %s is available (current %s); run `oym-conventions update`\n", tag, buildinfo.Version)
	}
}
