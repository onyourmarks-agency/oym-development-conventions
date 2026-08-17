package cards

import (
	"os"
	"strings"
	"testing"
)

func requireRepoRoot(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("../../conventions"); err != nil {
		t.Fatalf("repo conventions directory not found: %v", err)
	}
}

func TestEveryConventionFileHasOneValidCard(t *testing.T) {
	requireRepoRoot(t)
	assets := os.DirFS("../..")
	stacks, err := Stacks(assets)
	if err != nil {
		t.Fatal(err)
	}
	if len(stacks) == 0 {
		t.Fatal("no stacks found")
	}
	for _, stack := range stacks {
		card, err := Card(assets, stack)
		if err != nil {
			t.Errorf("stack %s: %v", stack, err)
			continue
		}
		lineCount := strings.Count(strings.TrimRight(card, "\n"), "\n") + 1
		if lineCount > MaxCardLines {
			t.Errorf("stack %s: card is %d lines, budget is %d", stack, lineCount, MaxCardLines)
		}
	}
}

const fixtureContent = `# AGENTS.md

Project rules here.

<!-- oym-conventions:begin stack=shared mode=existing -->
<!-- oym-conventions:end stack=shared -->
`

func TestSyncFillsAndStaysIdempotent(t *testing.T) {
	assets := os.DirFS("../..")

	first, err := SyncContent(fixtureContent, assets, "v1.0.0", "2026-08-17")
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Changed) != 1 || first.Changed[0] != "shared" {
		t.Fatalf("expected shared to change, got %v", first.Changed)
	}
	if !strings.Contains(first.Content, "mode=existing rev=v1.0.0 synced=2026-08-17") {
		t.Fatalf("header not restamped:\n%s", first.Content)
	}
	if !strings.Contains(first.Content, "Project rules here.") {
		t.Fatal("content outside blocks was lost")
	}

	second, err := SyncContent(first.Content, assets, "v9.9.9", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Changed) != 0 {
		t.Fatalf("expected idempotent second run, got changes %v", second.Changed)
	}
	if second.Content != first.Content {
		t.Fatal("second run modified content")
	}
}

func TestSyncRepairsTamperedBlock(t *testing.T) {
	assets := os.DirFS("../..")
	first, err := SyncContent(fixtureContent, assets, "v1.0.0", "2026-08-17")
	if err != nil {
		t.Fatal(err)
	}
	tampered := strings.Replace(first.Content, "# OYM shared rules", "# tampered", 1)

	repaired, err := SyncContent(tampered, assets, "v1.0.1", "2026-09-01")
	if err != nil {
		t.Fatal(err)
	}
	if len(repaired.Changed) != 1 {
		t.Fatalf("expected repair, got %v", repaired.Changed)
	}
	if !strings.Contains(repaired.Content, "# OYM shared rules") {
		t.Fatal("tampered card was not repaired")
	}
}

func TestSyncPreservesMode(t *testing.T) {
	assets := os.DirFS("../..")
	content := strings.Replace(fixtureContent, "mode=existing", "mode=new", 1)
	result, err := SyncContent(content, assets, "v1.0.0", "2026-08-17")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Content, "mode=new rev=v1.0.0") {
		t.Fatalf("mode=new was not preserved:\n%s", result.Content)
	}
}

func TestSyncErrorsOnUnclosedBlock(t *testing.T) {
	assets := os.DirFS("../..")
	_, err := SyncContent("<!-- oym-conventions:begin stack=shared mode=existing -->\nno end", assets, "v1", "d")
	if err == nil {
		t.Fatal("expected error for unclosed block")
	}
}

func TestSyncErrorsOnUnknownStack(t *testing.T) {
	assets := os.DirFS("../..")
	content := "<!-- oym-conventions:begin stack=nope -->\n<!-- oym-conventions:end stack=nope -->\n"
	_, err := SyncContent(content, assets, "v1", "d")
	if err == nil {
		t.Fatal("expected error for unknown stack")
	}
}

func TestCardsContainNoCodeFences(t *testing.T) {
	// The fence tracker assumes rendered block bodies never open a fence; a card
	// containing ``` would break scanning of everything after it.
	requireRepoRoot(t)
	assets := os.DirFS("../..")
	stacks, err := Stacks(assets)
	if err != nil {
		t.Fatal(err)
	}
	for _, stack := range stacks {
		card, err := Card(assets, stack)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(card, "```") || strings.Contains(card, "~~~") {
			t.Errorf("stack %s: card contains a code fence; cards must stay fence-free", stack)
		}
	}
}

func TestSyncIgnoresMarkersInsideCodeFences(t *testing.T) {
	assets := os.DirFS("../..")
	content := "# Docs\n\n```\n<!-- oym-conventions:begin stack=php mode=existing -->\n<!-- oym-conventions:end stack=php -->\n```\n\n" + fixtureContent
	result, err := SyncContent(content, assets, "v1", "2026-08-17")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changed) != 1 || result.Changed[0] != "shared" {
		t.Fatalf("only the live shared block should change, got %v", result.Changed)
	}
	if strings.Contains(result.Content, "```\n<!-- oym-conventions:begin stack=php mode=existing -->\n# OYM") {
		t.Fatal("card was injected into the fenced example")
	}

	present := PresentStacks(content)
	if len(present) != 1 || present[0] != "shared" {
		t.Fatalf("fenced markers must not count as present, got %v", present)
	}
	appended, added := AppendBlocks(content, []string{"php"}, "existing")
	if len(added) != 1 {
		t.Fatalf("php should be appendable despite the fenced example, got %v", added)
	}
	if !strings.Contains(appended, "\n<!-- oym-conventions:begin stack=php mode=existing -->\n<!-- oym-conventions:end stack=php -->\n") {
		t.Fatal("real php block missing after append")
	}
}

func TestSyncHandlesCRLFFiles(t *testing.T) {
	assets := os.DirFS("../..")
	synced, err := SyncContent(fixtureContent, assets, "v1", "2026-08-17")
	if err != nil {
		t.Fatal(err)
	}
	crlf := strings.ReplaceAll(synced.Content, "\n", "\r\n")

	unchanged, err := SyncContent(crlf, assets, "v2", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if len(unchanged.Changed) != 0 {
		t.Fatalf("CRLF file with current cards must not be stale, got %v", unchanged.Changed)
	}

	tampered := strings.Replace(crlf, "# OYM shared rules", "# tampered", 1)
	repaired, err := SyncContent(tampered, assets, "v2", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if len(repaired.Changed) != 1 {
		t.Fatalf("tampered CRLF block should be repaired, got %v", repaired.Changed)
	}
	if strings.Count(repaired.Content, "\n") != strings.Count(repaired.Content, "\r\n") {
		t.Fatal("repair must keep the file uniformly CRLF")
	}
}

func TestAppendBlocksDeduplicatesWithinRequest(t *testing.T) {
	out, added := AppendBlocks(fixtureContent, []string{"php", "php"}, "existing")
	if len(added) != 1 {
		t.Fatalf("duplicate request must append once, got %v", added)
	}
	if strings.Count(out, "begin stack=php") != 1 {
		t.Fatal("php block duplicated")
	}
}

func TestAppendBlocksSkipsPresent(t *testing.T) {
	out, added := AppendBlocks(fixtureContent, []string{"shared", "php"}, "existing")
	if len(added) != 1 || added[0] != "php" {
		t.Fatalf("expected only php appended, got %v", added)
	}
	if strings.Count(out, "begin stack=shared") != 1 {
		t.Fatal("shared block duplicated")
	}
	if !strings.Contains(out, "begin stack=php mode=existing") {
		t.Fatal("php block missing")
	}
}
