// Package project wires a product repository: AGENTS.md marker blocks, CLAUDE.md,
// and quality-gate template files.
package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/onyourmarks-agency/oym-development-conventions/internal/cards"
)

const (
	agentsFilename = "AGENTS.md"
	claudeFilename = "CLAUDE.md"
	orgHeading     = "## Organization conventions"

	agentsTemplatePath = "templates/AGENTS.template.md"
	claudeTemplatePath = "templates/CLAUDE.template.md"
)

type Options struct {
	Dir    string
	Assets fs.FS
	// Stacks to ensure a block exists for; empty means "sync existing blocks only".
	Stacks []string
	// Mode applied to newly appended blocks; existing blocks keep their own.
	Mode string
	Rev  string
	Date string
	// CheckOnly reports drift without writing.
	CheckOnly bool
}

type Result struct {
	AgentsPath    string
	AgentsCreated bool
	Appended      []string
	Updated       []string
}

func (r Result) UpToDate() bool {
	return !r.AgentsCreated && len(r.Appended) == 0 && len(r.Updated) == 0
}

// Sync ensures AGENTS.md exists, carries a block per requested stack, and has
// every block filled with the current card.
func Sync(opts Options) (Result, error) {
	if opts.Mode == "" {
		opts.Mode = cards.DefaultMode
	}
	result := Result{AgentsPath: filepath.Join(opts.Dir, agentsFilename)}

	if opts.Mode != "existing" && opts.Mode != "new" {
		return result, fmt.Errorf("invalid mode %q: use existing or new", opts.Mode)
	}
	if len(opts.Stacks) > 0 {
		available, err := cards.Stacks(opts.Assets)
		if err != nil {
			return result, err
		}
		for _, stack := range opts.Stacks {
			if !cards.ValidStackName(stack) || !slices.Contains(available, stack) {
				return result, fmt.Errorf("unknown stack %q (available: %v)", stack, available)
			}
		}
	}

	original, err := os.ReadFile(result.AgentsPath)
	content := string(original)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if len(opts.Stacks) == 0 {
			return result, fmt.Errorf("%s does not exist; run `oym-conventions init`", result.AgentsPath)
		}
		if opts.CheckOnly {
			return result, fmt.Errorf("%s does not exist", result.AgentsPath)
		}
		template, readErr := fs.ReadFile(opts.Assets, agentsTemplatePath)
		if readErr != nil {
			return result, fmt.Errorf("reading embedded AGENTS template: %w", readErr)
		}
		content = string(template)
		result.AgentsCreated = true
	case err != nil:
		return result, err
	}

	if len(opts.Stacks) > 0 {
		if len(cards.PresentStacks(content)) == 0 && !strings.Contains(content, orgHeading) {
			content = strings.TrimRight(content, "\n") + "\n\n" + orgHeading + "\n\n" +
				"Rules below are vendored from oym-development-conventions. Project-specific rules above override them. Do not edit inside the marker blocks.\n"
		}
		content, result.Appended = cards.AppendBlocks(content, opts.Stacks, opts.Mode)
	}

	synced, err := cards.SyncContent(content, opts.Assets, opts.Rev, opts.Date)
	if err != nil {
		return result, err
	}
	result.Updated = synced.Changed

	if opts.CheckOnly {
		return result, nil
	}
	if result.UpToDate() {
		return result, nil
	}
	if err := os.WriteFile(result.AgentsPath, []byte(synced.Content), 0o644); err != nil {
		return result, err
	}
	return result, nil
}

// EnsureClaudeMd writes CLAUDE.md from the embedded template when absent.
func EnsureClaudeMd(dir string, assets fs.FS) (bool, error) {
	path := filepath.Join(dir, claudeFilename)
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	template, err := fs.ReadFile(assets, claudeTemplatePath)
	if err != nil {
		return false, fmt.Errorf("reading embedded CLAUDE template: %w", err)
	}
	return true, os.WriteFile(path, template, 0o644)
}

type TemplateItem struct {
	Label  string
	Source string
	Dest   string
}

// TemplateItems returns the quality-gate templates that apply to the stack set.
func TemplateItems(stacks []string, mode string) []TemplateItem {
	items := []TemplateItem{}
	if slices.Contains(stacks, "php") {
		phpstan := "templates/php/phpstan.neon.dist"
		if mode == cards.DefaultMode {
			phpstan = "templates/php/phpstan-existing.neon.dist"
		}
		items = append(items,
			TemplateItem{Label: "phpcs.xml (OYM ruleset)", Source: "templates/php/phpcs.xml", Dest: "phpcs.xml"},
			TemplateItem{Label: "phpstan.neon.dist (level max)", Source: phpstan, Dest: "phpstan.neon.dist"},
			TemplateItem{Label: "phpunit.xml.dist (unit only)", Source: "templates/php/phpunit.xml.dist", Dest: "phpunit.xml.dist"},
		)
	}
	if slices.Contains(stacks, "tool-vite") {
		items = append(items, TemplateItem{Label: "vitest.config.ts (Vite islands)", Source: "templates/svelte/vitest.config.vite-craft.ts", Dest: "vitest.config.ts"})
	}
	if slices.Contains(stacks, "tool-sveltekit") {
		items = append(items, TemplateItem{Label: "vitest.config.ts (SvelteKit)", Source: "templates/svelte/vitest.config.sveltekit.ts", Dest: "vitest.config.ts"})
	}
	if slices.Contains(stacks, "framework-nestjs") {
		items = append(items,
			TemplateItem{Label: "tsconfig strictness delta", Source: "templates/nestjs/tsconfig.strict.delta.jsonc", Dest: "tsconfig.strict.delta.jsonc"},
			TemplateItem{Label: "eslint any-rules snippet", Source: "templates/nestjs/eslint.any-rules.snippet.mjs", Dest: "eslint.any-rules.snippet.mjs"},
		)
	}
	return items
}

// CopyTemplates copies the given items into dir, never overwriting existing files.
// It returns the destinations written and the destinations skipped.
func CopyTemplates(dir string, assets fs.FS, items []TemplateItem) (written []string, skipped []string, err error) {
	for _, item := range items {
		dest := filepath.Join(dir, item.Dest)
		if _, statErr := os.Stat(dest); statErr == nil {
			skipped = append(skipped, item.Dest)
			continue
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return written, skipped, statErr
		}
		data, readErr := fs.ReadFile(assets, item.Source)
		if readErr != nil {
			return written, skipped, fmt.Errorf("reading embedded template %s: %w", item.Source, readErr)
		}
		if writeErr := os.WriteFile(dest, data, 0o644); writeErr != nil {
			return written, skipped, writeErr
		}
		written = append(written, item.Dest)
	}
	return written, skipped, nil
}
