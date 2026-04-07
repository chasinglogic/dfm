package llm

import (
	"strings"
	"testing"
)

func TestBuildCommitMessagePromptUsesDefaultWhenUnset(t *testing.T) {
	diff := "diff --git a/file b/file"
	prompt := buildCommitMessagePrompt(diff, "")

	if !strings.Contains(prompt, "configuration-only diffs") {
		t.Fatalf("prompt should include default instructions")
	}

	if !strings.Contains(prompt, diff) {
		t.Fatalf("prompt should include diff")
	}
}

func TestBuildCommitMessagePromptUsesCustomTemplate(t *testing.T) {
	diff := "diff --git a/file b/file"
	template := "Custom instructions"
	prompt := buildCommitMessagePrompt(diff, template)

	if strings.Contains(prompt, "configuration-only diffs") {
		t.Fatalf("prompt should not include default instructions when custom template is set")
	}

	if !strings.Contains(prompt, "Custom instructions") {
		t.Fatalf("prompt should include custom instructions")
	}

	if !strings.Contains(prompt, "Diff:\n") {
		t.Fatalf("prompt should include the diff heading")
	}

	if !strings.Contains(prompt, diff) {
		t.Fatalf("prompt should include diff")
	}
}

func TestBuildCommitMessagePromptDoesNotInterpretPercentVerbs(t *testing.T) {
	diff := "diff --git a/file b/file"
	template := "Use this format literal %s and summarize changes"
	prompt := buildCommitMessagePrompt(diff, template)

	if !strings.Contains(prompt, "literal %s") {
		t.Fatalf("prompt should preserve literal percent verbs")
	}
}
