package llm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultCodexModel = "o4-mini"

// CodexProvider generates commit messages using the OpenAI Codex CLI.
// Requires the `codex` CLI to be installed and authenticated.
// This avoids the need for an OPENAI_API_KEY environment variable
// when authenticated via a ChatGPT subscription.
type CodexProvider struct {
	Model string
}

func (c *CodexProvider) GenerateCommitMessage(diff string, promptTemplate string) (string, error) {
	if _, err := exec.LookPath("codex"); err != nil {
		return "", fmt.Errorf("codex CLI not found in PATH: install it from https://github.com/openai/codex")
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	model := c.Model
	if model == "" {
		model = defaultCodexModel
	}

	// Write the output to a temp file since codex exec mixes progress
	// output on stderr and the response on stdout.
	tmpFile := filepath.Join(os.TempDir(), "dfm-codex-output.txt")
	defer os.Remove(tmpFile)

	cmd := exec.Command("codex", "exec", "--ephemeral", "-m", model, "-o", tmpFile, prompt)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("codex CLI failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run codex CLI: %w", err)
	}

	output, err := os.ReadFile(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to read codex output: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("empty response from codex CLI")
	}

	return result, nil
}
