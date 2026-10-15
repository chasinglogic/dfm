package llm

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewProviderPi(t *testing.T) {
	for _, model := range []string{"", "anthropic/claude-sonnet-4-6"} {
		p, err := NewProvider("pi", model)
		if err != nil {
			t.Fatal(err)
		}
		pp, ok := p.(*PiProvider)
		if !ok || pp.Model != model {
			t.Fatalf("expected PiProvider with model %q, got %#v", model, p)
		}
	}
}

func TestPiGenerateCommitMessage(t *testing.T) {
	for _, model := range []string{"", "anthropic/claude-sonnet-4-6"} {
		t.Run("model="+model, func(t *testing.T) {
			dir := t.TempDir()
			argsPath := filepath.Join(dir, "args")
			stdinPath := filepath.Join(dir, "stdin")
			t.Setenv("PI_TEST_ARGS", argsPath)
			t.Setenv("PI_TEST_STDIN", stdinPath)
			installFakePi(t, `printf '%s\n' "$@" > "$PI_TEST_ARGS"
while IFS= read -r line || [ -n "$line" ]; do
  printf '%s\n' "$line"
done > "$PI_TEST_STDIN"
printf '  Set editor default\n\n- Enable vim mode\n  '
`)
			p := &PiProvider{Model: model}
			diff := "-editor: emacs\n+editor: vim\n"
			message, err := p.GenerateCommitMessage(context.Background(), diff, "Write a commit message")
			if err != nil {
				t.Fatal(err)
			}
			if message != "Set editor default\n\n- Enable vim mode" {
				t.Fatalf("unexpected message: %q", message)
			}
			args, err := os.ReadFile(argsPath)
			if err != nil {
				t.Fatal(err)
			}
			wantArgs := []string{"--print", "--mode", "text", "--no-session", "--no-tools", "--no-extensions", "--no-skills", "--no-prompt-templates", "--no-context-files"}
			if model != "" {
				wantArgs = append(wantArgs, "--model", model)
			}
			if got := strings.Split(strings.TrimSpace(string(args)), "\n"); !reflect.DeepEqual(got, wantArgs) {
				t.Fatalf("args = %q, want %q", got, wantArgs)
			}
			input, err := os.ReadFile(stdinPath)
			if err != nil {
				t.Fatal(err)
			}
			if want := buildCommitMessagePrompt(diff, "Write a commit message"); string(input) != want {
				t.Fatalf("stdin = %q, want %q", input, want)
			}
		})
	}
}

func TestPiErrors(t *testing.T) {
	t.Run("missing executable", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		_, err := (&PiProvider{}).GenerateCommitMessage(context.Background(), "diff", "")
		if err == nil || !strings.Contains(err.Error(), "pi CLI not found in PATH") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	for _, tc := range []struct {
		name, script, want string
	}{
		{"empty response", "printf '  \\n'", "empty response from pi CLI"},
		{"exit failure", "printf 'authentication failed' >&2; exit 1", "authentication failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			installFakePi(t, tc.script)
			_, err := (&PiProvider{}).GenerateCommitMessage(context.Background(), "diff", "")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
	t.Run("canceled context", func(t *testing.T) {
		installFakePi(t, "exit 0")
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := (&PiProvider{}).GenerateCommitMessage(ctx, "diff", "")
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func installFakePi(t *testing.T, script string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pi"), []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
}
