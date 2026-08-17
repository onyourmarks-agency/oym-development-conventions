// Package cards extracts convention cards from the embedded conventions files and
// renders them into AGENTS.md marker blocks.
package cards

import (
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"sort"
	"strings"
)

const (
	conventionsDir  = "conventions"
	cardBeginFormat = "<!-- oym-card:begin stack=%s -->"
	cardEndFormat   = "<!-- oym-card:end stack=%s -->"
	blockEndFormat  = "<!-- oym-conventions:end stack=%s -->"

	DefaultMode  = "existing"
	MaxCardLines = 200
)

var (
	blockBeginPattern = regexp.MustCompile(`<!-- oym-conventions:begin stack=([a-z0-9-]+)((?: +[a-z]+=[^ >]+)*) *-->`)
	attributePattern  = regexp.MustCompile(`([a-z]+)=([^ >]+)`)
	stackNamePattern  = regexp.MustCompile(`^[a-z0-9-]+$`)
)

// ValidStackName reports whether a stack name can round-trip through the marker
// syntax. Anything else would render a begin marker the scanner cannot see.
func ValidStackName(stack string) bool {
	return stackNamePattern.MatchString(stack)
}

// Stacks lists every stack shipped in the embedded conventions directory.
func Stacks(assets fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(assets, conventionsDir)
	if err != nil {
		return nil, fmt.Errorf("reading embedded conventions: %w", err)
	}
	stacks := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}
		stacks = append(stacks, strings.TrimSuffix(name, ".md"))
	}
	sort.Strings(stacks)
	return stacks, nil
}

// Card returns the vendorable card body for a stack, normalized to end in exactly
// one newline and carry no leading newline.
func Card(assets fs.FS, stack string) (string, error) {
	raw, err := fs.ReadFile(assets, conventionsDir+"/"+stack+".md")
	if err != nil {
		return "", fmt.Errorf("unknown stack %q", stack)
	}
	content := string(raw)
	_, after, found := strings.Cut(content, fmt.Sprintf(cardBeginFormat, stack))
	if !found {
		return "", fmt.Errorf("conventions/%s.md has no card block", stack)
	}
	body, _, found := strings.Cut(after, fmt.Sprintf(cardEndFormat, stack))
	if !found {
		return "", fmt.Errorf("conventions/%s.md card block is not closed", stack)
	}
	return strings.Trim(body, "\n") + "\n", nil
}

type span struct {
	start int
	end   int
}

// fencedSpans returns the byte ranges of fenced code blocks (``` or ~~~), so
// marker syntax quoted in documentation is never treated as a live block.
func fencedSpans(content string) []span {
	spans := []span{}
	offset := 0
	inFence := false
	fenceStart := 0
	var fenceChar byte
	fenceLen := 0
	for _, line := range strings.SplitAfter(content, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		char, count := fenceHead(trimmed)
		if !inFence {
			if count >= 3 {
				inFence = true
				fenceChar = char
				fenceLen = count
				fenceStart = offset
			}
		} else if char == fenceChar && count >= fenceLen && strings.TrimRight(trimmed[count:], " \t\r\n") == "" {
			inFence = false
			spans = append(spans, span{fenceStart, offset + len(line)})
		}
		offset += len(line)
	}
	if inFence {
		spans = append(spans, span{fenceStart, len(content)})
	}
	return spans
}

func fenceHead(line string) (byte, int) {
	if line == "" || (line[0] != '`' && line[0] != '~') {
		return 0, 0
	}
	char := line[0]
	count := 0
	for count < len(line) && line[count] == char {
		count++
	}
	return char, count
}

func insideSpans(spans []span, position int) bool {
	for _, region := range spans {
		if position >= region.start && position < region.end {
			return true
		}
	}
	return false
}

// PresentStacks returns the stacks that have a live marker block in content, in
// order, ignoring markers inside fenced code blocks.
func PresentStacks(content string) []string {
	fences := fencedSpans(content)
	stacks := []string{}
	for _, match := range blockBeginPattern.FindAllStringSubmatchIndex(content, -1) {
		if insideSpans(fences, match[0]) {
			continue
		}
		stacks = append(stacks, content[match[2]:match[3]])
	}
	return stacks
}

// AppendBlocks appends an empty marker block for every requested stack that has no
// live block yet and returns the new content plus the appended stacks. Duplicates
// within one request are appended once.
func AppendBlocks(content string, stacks []string, mode string) (string, []string) {
	present := PresentStacks(content)
	added := []string{}
	out := strings.TrimRight(content, "\n")
	for _, stack := range stacks {
		if slices.Contains(present, stack) {
			continue
		}
		present = append(present, stack)
		added = append(added, stack)
		out += fmt.Sprintf(
			"\n\n<!-- oym-conventions:begin stack=%s mode=%s -->\n"+blockEndFormat,
			stack, mode, stack,
		)
	}
	return out + "\n", added
}

type SyncResult struct {
	Content string
	Changed []string
}

// SyncContent replaces the body of every live marker block with the current card
// for its stack. The mode attribute is preserved; rev and synced are restamped
// only on blocks whose body actually changed, so the operation is idempotent.
// CRLF files are compared EOL-insensitively and rewritten with their own EOL.
func SyncContent(content string, assets fs.FS, rev string, date string) (SyncResult, error) {
	fences := fencedSpans(content)
	eol := "\n"
	if strings.Contains(content, "\r\n") {
		eol = "\r\n"
	}

	var out strings.Builder
	changed := []string{}
	cursor := 0
	for {
		loc := blockBeginPattern.FindStringSubmatchIndex(content[cursor:])
		if loc == nil {
			out.WriteString(content[cursor:])
			break
		}
		blockStart := cursor + loc[0]
		headerEnd := cursor + loc[1]
		if insideSpans(fences, blockStart) {
			out.WriteString(content[cursor:headerEnd])
			cursor = headerEnd
			continue
		}
		stack := content[cursor+loc[2] : cursor+loc[3]]
		attributes := parseAttributes(content[cursor+loc[4] : cursor+loc[5]])

		endMarkerStart, endMarkerLength, found := findEndMarker(content, headerEnd, stack, fences)
		if !found {
			return SyncResult{}, fmt.Errorf("marker block for stack %q is not closed", stack)
		}
		blockEnd := endMarkerStart + endMarkerLength

		card, err := Card(assets, stack)
		if err != nil {
			return SyncResult{}, err
		}

		out.WriteString(content[cursor:blockStart])
		body := content[headerEnd:endMarkerStart]
		normalizedBody := strings.ReplaceAll(body, "\r\n", "\n")
		if strings.Trim(normalizedBody, "\n")+"\n" == card {
			out.WriteString(content[blockStart:blockEnd])
		} else {
			changed = append(changed, stack)
			mode := attributes["mode"]
			if mode == "" {
				mode = DefaultMode
			}
			header := fmt.Sprintf("<!-- oym-conventions:begin stack=%s mode=%s rev=%s synced=%s -->", stack, mode, rev, date)
			out.WriteString(header + eol + strings.ReplaceAll(card, "\n", eol) + fmt.Sprintf(blockEndFormat, stack))
		}
		cursor = blockEnd
	}
	return SyncResult{Content: out.String(), Changed: changed}, nil
}

func findEndMarker(content string, from int, stack string, fences []span) (start int, length int, found bool) {
	endMarker := fmt.Sprintf(blockEndFormat, stack)
	searchFrom := from
	for {
		offset := strings.Index(content[searchFrom:], endMarker)
		if offset < 0 {
			return 0, 0, false
		}
		position := searchFrom + offset
		if !insideSpans(fences, position) {
			return position, len(endMarker), true
		}
		searchFrom = position + len(endMarker)
	}
}

func parseAttributes(text string) map[string]string {
	attributes := map[string]string{}
	for _, match := range attributePattern.FindAllStringSubmatch(text, -1) {
		attributes[match[1]] = match[2]
	}
	return attributes
}
