// Package detect derives the applicable convention stacks from a project's
// composer.json, package.json, and directory layout.
package detect

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type composerFile struct {
	Type    string            `json:"type"`
	Require map[string]string `json:"require"`
}

type packageFile struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// Stacks returns the convention stacks detected in dir, always starting with
// "shared". Unknown project types yield just the layers that could be proven.
func Stacks(dir string) []string {
	stacks := []string{"shared"}
	stacks = append(stacks, phpStacks(dir)...)
	stacks = append(stacks, javascriptStacks(dir)...)
	return stacks
}

func phpStacks(dir string) []string {
	composer, err := readJSON[composerFile](filepath.Join(dir, "composer.json"))
	if err != nil {
		return nil
	}
	switch {
	case composer.Type == "craft-plugin":
		return []string{"php", "framework-craft-plugin"}
	case hasKey(composer.Require, "craftcms/cms"):
		stacks := []string{"php", "framework-craft"}
		if dirExists(filepath.Join(dir, "api-simple")) {
			stacks = append(stacks, "framework-api-simple")
		}
		return stacks
	case hasKey(composer.Require, "symfony/framework-bundle"):
		return []string{"php", "framework-symfony"}
	case len(composer.Require) > 0:
		return []string{"php"}
	}
	return nil
}

func javascriptStacks(dir string) []string {
	pkg, err := readJSON[packageFile](filepath.Join(dir, "package.json"))
	if err != nil {
		return nil
	}
	dependencies := map[string]string{}
	for name, version := range pkg.Dependencies {
		dependencies[name] = version
	}
	for name, version := range pkg.DevDependencies {
		dependencies[name] = version
	}
	switch {
	case hasKey(dependencies, "@nestjs/core"):
		return []string{"javascript", "framework-nestjs"}
	case hasKey(dependencies, "@sveltejs/kit"):
		return []string{"javascript", "framework-svelte", "tool-sveltekit"}
	case hasKey(dependencies, "svelte"):
		tool := "tool-vite"
		if dirExists(filepath.Join(dir, "webpack")) {
			tool = "tool-webpack"
		}
		return []string{"javascript", "framework-svelte", tool}
	case len(dependencies) > 0:
		return []string{"javascript"}
	}
	return nil
}

func readJSON[T any](path string) (T, error) {
	var parsed T
	data, err := os.ReadFile(path)
	if err != nil {
		return parsed, err
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return parsed, err
	}
	return parsed, nil
}

func hasKey(values map[string]string, key string) bool {
	_, ok := values[key]
	return ok
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
