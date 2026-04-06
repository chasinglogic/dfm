package cmd

import (
	"regexp"
	"testing"
)

func TestLinkAsDirMatchPatternMatchesProvidedDirectoryOnly(t *testing.T) {
	pattern := linkAsDirMatchPattern(".agents")
	rgx := regexp.MustCompile(pattern)

	if !rgx.MatchString("/tmp/repo/.agents") {
		t.Fatalf("pattern %q should match directory root", pattern)
	}

	if rgx.MatchString("/tmp/repo/.agents/skills/test.md") {
		t.Fatalf("pattern %q should not match directory children", pattern)
	}

	if rgx.MatchString("/tmp/repo/.agentsx") {
		t.Fatalf("pattern %q should not match sibling paths", pattern)
	}
}

func TestLinkAsDirMatchPatternEscapesRegexChars(t *testing.T) {
	pattern := linkAsDirMatchPattern(".config/nvim/snippets")
	rgx := regexp.MustCompile(pattern)

	if !rgx.MatchString("/tmp/repo/.config/nvim/snippets") {
		t.Fatalf("pattern %q should match literal dotted directory", pattern)
	}

	if rgx.MatchString("/tmp/repo/xconfig/nvim/snippets") {
		t.Fatalf("pattern %q should treat dots as literals", pattern)
	}
}
