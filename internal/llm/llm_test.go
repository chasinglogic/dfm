package llm

import (
	"strings"
	"testing"
)

func TestBuildCommitMessagePromptUsesDefaultWhenUnset(t *testing.T) {
	diff := "diff --git a/file b/file"
	prompt := buildCommitMessagePrompt(diff, "")

	if !strings.Contains(prompt, "Generate a concise") {
		t.Fatalf("prompt should include default instructions")
	}

	if !strings.Contains(prompt, diff) {
		t.Fatalf("prompt should include diff")
	}
}

func TestBuildCommitMessagePromptUsesCustomTemplate(t *testing.T) {
	diff := "diff --git a/file b/file"
	template := "Custom instructions\n\nDiff:\n%s"
	prompt := buildCommitMessagePrompt(diff, template)

	if strings.Contains(prompt, "Generate a concise") {
		t.Fatalf("prompt should not include default instructions when custom template is set")
	}

	if !strings.Contains(prompt, "Custom instructions") {
		t.Fatalf("prompt should include custom instructions")
	}

	if !strings.Contains(prompt, diff) {
		t.Fatalf("prompt should include diff")
	}
}
