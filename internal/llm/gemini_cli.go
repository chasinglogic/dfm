package llm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/chasinglogic/dfm/internal/logger"
)

const defaultGeminiCLIModel = "gemini-2.5-flash"

// GeminiCLIProvider generates commit messages using the Gemini CLI.
// Requires the `gemini` CLI to be installed and authenticated.
// This avoids the need for a GEMINI_API_KEY environment variable.
type GeminiCLIProvider struct {
	Model string
}

func (g *GeminiCLIProvider) GenerateCommitMessage(ctx context.Context, diff string, promptTemplate string) (string, error) {
	if _, err := exec.LookPath("gemini"); err != nil {
		return "", fmt.Errorf("gemini CLI not found in PATH: install it from https://github.com/google-gemini/gemini-cli")
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	model := g.Model
	if model == "" {
		model = defaultGeminiCLIModel
	}

	logger.Debug().Str("provider", "gemini-cli").Str("model", model).Int("diffBytes", len(diff)).Msg("running gemini cli commit message request")

	started := time.Now()
	cmd := exec.CommandContext(ctx, "gemini", "-p", prompt, "-m", model)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	logger.Debug().Str("provider", "gemini-cli").Msg("starting gemini cli subprocess")
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("gemini CLI failed: %s", strings.TrimSpace(stderr.String()+" "+string(exitErr.Stderr)))
		}
		if ctx.Err() != nil {
			return "", fmt.Errorf("gemini CLI timed out after %s: %w", time.Since(started).Truncate(time.Second), ctx.Err())
		}
		return "", fmt.Errorf("failed to run gemini CLI: %w", err)
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", fmt.Errorf("empty response from gemini CLI")
	}

	logger.Debug().Str("provider", "gemini-cli").Dur("elapsed", time.Since(started)).Int("messageBytes", len(result)).Msg("finished gemini cli commit message request")

	return result, nil
}
