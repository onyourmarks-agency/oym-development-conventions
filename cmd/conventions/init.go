package conventions

import (
	"fmt"
	"io/fs"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/buildinfo"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/cards"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/detect"
	"github.com/onyourmarks-agency/oym-development-conventions/internal/project"
)

type initOptions struct {
	target        string
	stacks        string
	mode          string
	yes           bool
	withTemplates bool
	noTemplates   bool
}

func NewInitCmd(assets fs.FS) *cobra.Command {
	o := &initOptions{}
	c := &cobra.Command{
		Use:   "init",
		Short: "Wire a project: detect stacks, write AGENTS.md blocks, CLAUDE.md, and quality-gate templates.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runInit(assets, o)
		},
	}
	c.Flags().StringVar(&o.target, "target", ".", "project directory")
	c.Flags().StringVar(&o.stacks, "stacks", "", "comma-separated stacks (default: detected)")
	c.Flags().StringVar(&o.mode, "mode", "", "existing or new (default: existing)")
	c.Flags().BoolVarP(&o.yes, "yes", "y", false, "non-interactive: accept detected stacks and defaults")
	c.Flags().BoolVar(&o.withTemplates, "with-templates", false, "copy all applicable quality-gate templates")
	c.Flags().BoolVar(&o.noTemplates, "no-templates", false, "skip quality-gate templates")
	c.MarkFlagsMutuallyExclusive("with-templates", "no-templates")
	return c
}

func runInit(assets fs.FS, o *initOptions) error {
	available, err := cards.Stacks(assets)
	if err != nil {
		return err
	}

	selected := splitStacks(o.stacks)
	if len(selected) == 0 {
		selected = detect.Stacks(o.target)
	}
	for _, stack := range selected {
		if !slices.Contains(available, stack) {
			return fmt.Errorf("unknown stack %q (available: %v)", stack, available)
		}
	}
	mode := o.mode
	if mode == "" {
		mode = cards.DefaultMode
	}
	if mode != "existing" && mode != "new" {
		return fmt.Errorf("invalid mode %q: use existing or new", mode)
	}

	interactive := !o.yes && stdinIsTerminal()
	if interactive {
		options := make([]huh.Option[string], 0, len(available))
		for _, stack := range available {
			options = append(options, huh.NewOption(stack, stack).Selected(slices.Contains(selected, stack)))
		}
		form := huh.NewForm(huh.NewGroup(
			huh.NewNote().
				Title("OYM conventions").
				Description("Wires AGENTS.md, CLAUDE.md, and quality-gate templates.\nDetected stacks are preselected."),
			huh.NewMultiSelect[string]().
				Title("Convention stacks").
				Options(options...).
				Validate(func(chosen []string) error {
					if len(chosen) == 0 {
						return fmt.Errorf("select at least one stack")
					}
					return nil
				}).
				Value(&selected),
			huh.NewSelect[string]().
				Title("Project mode").
				Description("existing: mirror local code, tooling ratchets in gently.\nnew: full standard from day one.").
				Options(huh.NewOption("existing", "existing"), huh.NewOption("new", "new")).
				Value(&mode),
		))
		if err := form.Run(); err != nil {
			return err
		}
	}

	result, err := project.Sync(project.Options{
		Dir:    o.target,
		Assets: assets,
		Stacks: selected,
		Mode:   mode,
		Rev:    buildinfo.Version,
		Date:   today(),
	})
	if err != nil {
		return err
	}
	claudeCreated, err := project.EnsureClaudeMd(o.target, assets)
	if err != nil {
		return err
	}
	report(result, claudeCreated)

	items := project.TemplateItems(selected, mode)
	copyTemplates := o.withTemplates || (mode == "new" && !o.noTemplates)
	if interactive && !o.withTemplates && !o.noTemplates && len(items) > 0 {
		chosen := []project.TemplateItem{}
		options := make([]huh.Option[project.TemplateItem], 0, len(items))
		for _, item := range items {
			options = append(options, huh.NewOption(item.Label, item).Selected(mode == "new"))
		}
		form := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[project.TemplateItem]().
				Title("Quality-gate templates to copy").
				Options(options...).
				Value(&chosen),
		))
		if err := form.Run(); err != nil {
			return err
		}
		items = chosen
		copyTemplates = len(items) > 0
	}

	if copyTemplates && len(items) > 0 {
		written, skipped, err := project.CopyTemplates(o.target, assets, items)
		if err != nil {
			return err
		}
		for _, dest := range written {
			fmt.Printf("  wrote %s\n", dest)
		}
		for _, dest := range skipped {
			fmt.Printf("  kept existing %s\n", dest)
		}
	}
	fmt.Println("done; review AGENTS.md and fill in the project sections")
	return nil
}

func report(result project.Result, claudeCreated bool) {
	if result.AgentsCreated {
		fmt.Printf("  created %s\n", result.AgentsPath)
	}
	if len(result.Appended) > 0 {
		fmt.Printf("  added blocks: %v\n", result.Appended)
	}
	if len(result.Updated) > 0 {
		fmt.Printf("  filled blocks: %v\n", result.Updated)
	}
	if claudeCreated {
		fmt.Println("  created CLAUDE.md")
	}
	if result.UpToDate() && !claudeCreated {
		fmt.Println("  already up to date")
	}
}
