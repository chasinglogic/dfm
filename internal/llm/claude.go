package llm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/chasinglogic/dfm/internal/logger"
)

const defaultClaudeModel = "sonnet"

// ClaudeProvider generates commit messages using the Claude CLI.
// Requires the `claude` CLI to be installed and authenticated.
type ClaudeProvider struct {
	Model string
}

func (c *ClaudeProvider) GenerateCommitMessage(ctx context.Context, diff string, promptTemplate string) (string, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return "", fmt.Errorf("claude CLI not found in PATH: install it from https://docs.anthropic.com/en/docs/claude-cli")
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	model := c.Model
	if model == "" {
		model = defaultClaudeModel
	}

	logger.Debug().Str("provider", "claude").Str("model", model).Int("diffBytes", len(diff)).Msg("running claude commit message request")

	started := time.Now()
	cmd := exec.CommandContext(ctx, "claude", "--print", "--model", model, prompt)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude CLI failed: %s", string(exitErr.Stderr))
		}
		if ctx.Err() != nil {
			return "", fmt.Errorf("claude CLI timed out after %s: %w", time.Since(started).Truncate(time.Second), ctx.Err())
		}
		return "", fmt.Errorf("failed to run claude CLI: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("empty response from claude CLI")
	}

	logger.Debug().Str("provider", "claude").Dur("elapsed", time.Since(started)).Int("messageBytes", len(result)).Msg("finished claude commit message request")

	return result, nil
}
