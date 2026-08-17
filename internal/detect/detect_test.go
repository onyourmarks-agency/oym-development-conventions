package detect

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeFixture(t *testing.T, dir string, name string, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestStacks(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(t *testing.T, dir string)
		expected []string
	}{
		{
			name:     "empty repo",
			setup:    func(_ *testing.T, _ string) {},
			expected: []string{"shared"},
		},
		{
			name: "craft site with api-simple and vite frontend",
			setup: func(t *testing.T, dir string) {
				writeFixture(t, dir, "composer.json", `{"require":{"craftcms/cms":"^5.0"}}`)
				writeFixture(t, dir, "package.json", `{"devDependencies":{"svelte":"^5.0.0","vite":"^7.0.0"}}`)
				if err := os.Mkdir(filepath.Join(dir, "api-simple"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			expected: []string{"shared", "php", "framework-craft", "framework-api-simple", "javascript", "framework-svelte", "tool-vite"},
		},
		{
			name: "craft site with webpack frontend",
			setup: func(t *testing.T, dir string) {
				writeFixture(t, dir, "composer.json", `{"require":{"craftcms/cms":"^4.0"}}`)
				writeFixture(t, dir, "package.json", `{"devDependencies":{"svelte":"^5.0.0"}}`)
				if err := os.Mkdir(filepath.Join(dir, "webpack"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			expected: []string{"shared", "php", "framework-craft", "javascript", "framework-svelte", "tool-webpack"},
		},
		{
			name: "craft plugin",
			setup: func(t *testing.T, dir string) {
				writeFixture(t, dir, "composer.json", `{"type":"craft-plugin","require":{"craftcms/cms":"^5.0"}}`)
			},
			expected: []string{"shared", "php", "framework-craft-plugin"},
		},
		{
			name: "symfony api",
			setup: func(t *testing.T, dir string) {
				writeFixture(t, dir, "composer.json", `{"require":{"symfony/framework-bundle":"^7.0"}}`)
			},
			expected: []string{"shared", "php", "framework-symfony"},
		},
		{
			name: "sveltekit app",
			setup: func(t *testing.T, dir string) {
				writeFixture(t, dir, "package.json", `{"devDependencies":{"@sveltejs/kit":"^2.0.0","svelte":"^5.0.0"}}`)
			},
			expected: []string{"shared", "javascript", "framework-svelte", "tool-sveltekit"},
		},
		{
			name: "nestjs service",
			setup: func(t *testing.T, dir string) {
				writeFixture(t, dir, "package.json", `{"dependencies":{"@nestjs/core":"^11.0.0"}}`)
			},
			expected: []string{"shared", "javascript", "framework-nestjs"},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			dir := t.TempDir()
			testCase.setup(t, dir)
			stacks := Stacks(dir)
			if !reflect.DeepEqual(stacks, testCase.expected) {
				t.Fatalf("expected %v, got %v", testCase.expected, stacks)
			}
		})
	}
}
