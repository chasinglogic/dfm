package profiles

import (
	"context"
	"fmt"
	"time"

	"github.com/chasinglogic/dfm/internal/llm"
	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/chasinglogic/dfm/internal/utils"
	"github.com/chzyer/readline"
)

const llmCommitMessageTimeout = 2 * time.Minute

func commitMessageFromLLM(location string, provider string, model string, prompt string) (string, error) {
	started := time.Now()
	logger.Debug().
		Str("location", location).
		Str("provider", provider).
		Str("model", model).
		Dur("timeout", llmCommitMessageTimeout).
		Msg("preparing LLM commit message")

	if err := utils.RunIn(location, "git", "add", "--all"); err != nil {
		return "", err
	}
	logger.Debug().Str("location", location).Msg("staged changes for LLM diff")

	diff, err := utils.RunInOutput(location, "git", "diff", "--cached")
	if err == nil && diff != "" {
		logger.Debug().Str("location", location).Int("diffBytes", len(diff)).Msg("collected staged diff for LLM")

		ctx, cancel := context.WithTimeout(context.Background(), llmCommitMessageTimeout)
		defer cancel()

		msg, llmErr := llm.GenerateCommitMessage(ctx, diff, provider, model, prompt)
		if llmErr != nil {
			return "", fmt.Errorf("failed to generate commit message from LLM (you enabled LLM commit messages with %s provider): %w", provider, llmErr)
		}

		logger.Debug().
			Str("location", location).
			Dur("elapsed", time.Since(started)).
			Int("messageBytes", len(msg)).
			Msg("generated LLM commit message")

		return msg, nil
	}

	if err != nil {
		logger.Debug().Str("location", location).Err(err).Msg("failed to read staged diff for LLM commit message")
	}
	logger.Debug().Str("location", location).Msg("no staged diff for LLM commit message")

	return "", nil
}

func commitMessageFromPrompt(location string) (string, error) {
	if err := utils.RunIn(location, "git", "diff"); err != nil {
		return "", err
	}

	rl, err := readline.New("Commit message: ")
	if err != nil {
		panic(err)
	}

	commitMessage, err := rl.Readline()
	if err != nil {
		return "", err
	}

	return commitMessage, rl.Close()
}
