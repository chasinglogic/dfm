package llm

import (
	"fmt"
	"os/exec"
	"strings"
)

const defaultClaudeModel = "sonnet"

// ClaudeProvider generates commit messages using the Claude CLI.
// Requires the `claude` CLI to be installed and authenticated.
type ClaudeProvider struct {
	Model string
}

func (c *ClaudeProvider) GenerateCommitMessage(diff string, promptTemplate string) (string, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return "", fmt.Errorf("claude CLI not found in PATH: install it from https://docs.anthropic.com/en/docs/claude-cli")
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	model := c.Model
	if model == "" {
		model = defaultClaudeModel
	}

	cmd := exec.Command("claude", "--print", "--model", model, prompt)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude CLI failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude CLI: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("empty response from claude CLI")
	}

	return result, nil
}
