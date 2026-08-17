package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func repoAssets(t *testing.T) fs.FS {
	t.Helper()
	if _, err := os.Stat("../../conventions"); err != nil {
		t.Fatalf("repo conventions directory not found: %v", err)
	}
	return os.DirFS("../..")
}

func options(dir string, assets fs.FS, stacks []string, mode string) Options {
	return Options{
		Dir:    dir,
		Assets: assets,
		Stacks: stacks,
		Mode:   mode,
		Rev:    "v1.0.0",
		Date:   "2026-08-17",
	}
}

func TestLifecycleFromScratch(t *testing.T) {
	assets := repoAssets(t)
	dir := t.TempDir()
	stacks := []string{"shared", "php", "framework-craft"}

	created, err := Sync(options(dir, assets, stacks, "existing"))
	if err != nil {
		t.Fatal(err)
	}
	if !created.AgentsCreated {
		t.Fatal("AGENTS.md should be created")
	}
	content := readFile(t, filepath.Join(dir, "AGENTS.md"))
	for _, marker := range []string{"# OYM shared rules", "# OYM PHP rules", "# OYM Craft CMS site rules"} {
		if !strings.Contains(content, marker) {
			t.Fatalf("missing card %q", marker)
		}
	}

	second, err := Sync(options(dir, assets, nil, ""))
	if err != nil {
		t.Fatal(err)
	}
	if !second.UpToDate() {
		t.Fatalf("second sync should be up to date, got %+v", second)
	}
	if readFile(t, filepath.Join(dir, "AGENTS.md")) != content {
		t.Fatal("second sync modified the file")
	}

	tampered := strings.Replace(content, "# OYM PHP rules", "# tampered", 1)
	writeFile(t, filepath.Join(dir, "AGENTS.md"), tampered)

	check, err := Sync(Options{Dir: dir, Assets: assets, Rev: "v1", Date: "d", CheckOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(check.Updated) != 1 || check.Updated[0] != "php" {
		t.Fatalf("check should flag php drift, got %+v", check)
	}
	if readFile(t, filepath.Join(dir, "AGENTS.md")) != tampered {
		t.Fatal("check must not write")
	}

	repaired, err := Sync(options(dir, assets, nil, ""))
	if err != nil {
		t.Fatal(err)
	}
	if len(repaired.Updated) != 1 {
		t.Fatalf("sync should repair the tampered block, got %+v", repaired)
	}
	if !strings.Contains(readFile(t, filepath.Join(dir, "AGENTS.md")), "# OYM PHP rules") {
		t.Fatal("tampered card not repaired")
	}
}

func TestAppendToExistingAgentsFileWithoutMarkers(t *testing.T) {
	assets := repoAssets(t)
	dir := t.TempDir()
	existing := "# AGENTS.md\n\nProject documentation that must survive.\n"
	writeFile(t, filepath.Join(dir, "AGENTS.md"), existing)

	result, err := Sync(options(dir, assets, []string{"javascript", "framework-nestjs"}, "existing"))
	if err != nil {
		t.Fatal(err)
	}
	if result.AgentsCreated {
		t.Fatal("existing file must not be recreated")
	}
	content := readFile(t, filepath.Join(dir, "AGENTS.md"))
	if !strings.HasPrefix(content, existing[:20]) {
		t.Fatal("existing content lost")
	}
	if !strings.Contains(content, "## Organization conventions") {
		t.Fatal("organization heading missing")
	}
	if !strings.Contains(content, "# OYM NestJS rules") {
		t.Fatal("nestjs card missing")
	}
}

func TestSyncRejectsInvalidStackAndMode(t *testing.T) {
	assets := repoAssets(t)
	dir := t.TempDir()

	if _, err := Sync(options(dir, assets, []string{"PHP"}, "existing")); err == nil {
		t.Fatal("uppercase stack must be rejected")
	}
	if _, err := Sync(options(dir, assets, []string{"phpp"}, "existing")); err == nil {
		t.Fatal("unknown stack must be rejected")
	}
	if _, err := Sync(options(dir, assets, []string{"php"}, "brand new")); err == nil {
		t.Fatal("invalid mode must be rejected")
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatal("rejected input must not write any file")
	}
}

func TestCheckOnlyOnMissingFileFails(t *testing.T) {
	assets := repoAssets(t)
	dir := t.TempDir()
	_, err := Sync(Options{Dir: dir, Assets: assets, Stacks: []string{"shared"}, Rev: "v1", Date: "d", CheckOnly: true})
	if err == nil {
		t.Fatal("check on a missing AGENTS.md must fail")
	}
}

func TestEnsureClaudeMd(t *testing.T) {
	assets := repoAssets(t)
	dir := t.TempDir()

	created, err := EnsureClaudeMd(dir, assets)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Fatal("CLAUDE.md should be created")
	}
	if !strings.Contains(readFile(t, filepath.Join(dir, "CLAUDE.md")), "@AGENTS.md") {
		t.Fatal("CLAUDE.md missing @AGENTS.md import")
	}

	again, err := EnsureClaudeMd(dir, assets)
	if err != nil {
		t.Fatal(err)
	}
	if again {
		t.Fatal("existing CLAUDE.md must not be recreated")
	}
}

func TestCopyTemplatesNeverOverwrites(t *testing.T) {
	assets := repoAssets(t)
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "phpcs.xml"), "custom")

	items := TemplateItems([]string{"shared", "php"}, "new")
	written, skipped, err := CopyTemplates(dir, assets, items)
	if err != nil {
		t.Fatal(err)
	}
	if readFile(t, filepath.Join(dir, "phpcs.xml")) != "custom" {
		t.Fatal("existing phpcs.xml was overwritten")
	}
	if len(skipped) != 1 || skipped[0] != "phpcs.xml" {
		t.Fatalf("expected phpcs.xml skipped, got %v", skipped)
	}
	if len(written) != 2 {
		t.Fatalf("expected phpstan + phpunit written, got %v", written)
	}
	phpstan := readFile(t, filepath.Join(dir, "phpstan.neon.dist"))
	if strings.Contains(phpstan, "phpstan-baseline.neon") {
		t.Fatal("mode=new must use the no-baseline phpstan template")
	}
}

func TestTemplateItemsExistingModeUsesBaselineVariant(t *testing.T) {
	items := TemplateItems([]string{"php"}, "existing")
	found := false
	for _, item := range items {
		if item.Dest == "phpstan.neon.dist" && strings.Contains(item.Source, "phpstan-existing") {
			found = true
		}
	}
	if !found {
		t.Fatal("existing mode should select the baseline phpstan template")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
