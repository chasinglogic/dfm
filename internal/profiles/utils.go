package profiles

import (
	"fmt"

	"github.com/chasinglogic/dfm/internal/llm"
	"github.com/chasinglogic/dfm/internal/utils"
	"github.com/chzyer/readline"
)

func commitMessageFromLLM(location string, provider string) (string, error) {
	if err := utils.RunIn(location, "git", "add", "--all"); err != nil {
		return "", err
	}

	diff, err := utils.RunInOutput(location, "git", "diff", "--cached")
	if err == nil && diff != "" {
		msg, llmErr := llm.GenerateCommitMessage(diff, provider)
		if llmErr != nil {
			return "", fmt.Errorf("failed to generate commit message from LLM (you enabled LLM commit messages with %s provider): %w", provider, llmErr)
		}

		return msg, nil
	}

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
