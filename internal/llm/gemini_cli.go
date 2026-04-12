package llm

import (
	"fmt"
	"os/exec"
	"strings"
)

const defaultGeminiCLIModel = "gemini-2.5-flash"

// GeminiCLIProvider generates commit messages using the Gemini CLI.
// Requires the `gemini` CLI to be installed and authenticated.
// This avoids the need for a GEMINI_API_KEY environment variable.
type GeminiCLIProvider struct {
	Model string
}

func (g *GeminiCLIProvider) GenerateCommitMessage(diff string, promptTemplate string) (string, error) {
	if _, err := exec.LookPath("gemini"); err != nil {
		return "", fmt.Errorf("gemini CLI not found in PATH: install it from https://github.com/google-gemini/gemini-cli")
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	model := g.Model
	if model == "" {
		model = defaultGeminiCLIModel
	}

	cmd := exec.Command("gemini", "-p", prompt, "-m", model)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("gemini CLI failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run gemini CLI: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("empty response from gemini CLI")
	}

	return result, nil
}
