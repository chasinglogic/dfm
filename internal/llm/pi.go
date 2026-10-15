package llm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/chasinglogic/dfm/internal/logger"
)

// PiProvider generates commit messages using the authenticated pi CLI.
// An empty Model uses pi's configured default.
type PiProvider struct {
	Model string
}

func (p *PiProvider) GenerateCommitMessage(ctx context.Context, diff string, promptTemplate string) (string, error) {
	if _, err := exec.LookPath("pi"); err != nil {
		return "", fmt.Errorf("pi CLI not found in PATH: install it from https://pi.dev")
	}

	// Keep generation text-only and independent of coding-agent resources.
	args := []string{"--print", "--mode", "text", "--no-session", "--no-tools", "--no-extensions", "--no-skills", "--no-prompt-templates", "--no-context-files"}
	if p.Model != "" {
		args = append(args, "--model", p.Model)
	}

	logger.Debug().Str("provider", "pi").Str("model", p.Model).Int("diffBytes", len(diff)).Msg("running pi commit message request")
	started := time.Now()
	cmd := exec.CommandContext(ctx, "pi", args...)
	cmd.Stdin = strings.NewReader(buildCommitMessagePrompt(diff, promptTemplate))
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("pi CLI canceled after %s: %w", time.Since(started).Truncate(time.Second), ctx.Err())
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("pi CLI failed: %s: %w", strings.TrimSpace(string(exitErr.Stderr)), err)
		}
		return "", fmt.Errorf("failed to run pi CLI: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "", fmt.Errorf("empty response from pi CLI")
	}

	logger.Debug().Str("provider", "pi").Dur("elapsed", time.Since(started)).Int("messageBytes", len(result)).Msg("finished pi commit message request")
	return result, nil
}
