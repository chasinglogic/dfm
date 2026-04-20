package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/openai/openai-go/v3"
)

const defaultOpenAIModel = "gpt-4.1-mini"

// OpenAIProvider generates commit messages using the OpenAI API.
// Requires the OPENAI_API_KEY environment variable to be set.
type OpenAIProvider struct {
	Model string
}

func (o *OpenAIProvider) GenerateCommitMessage(ctx context.Context, diff string, promptTemplate string) (string, error) {
	logger.Debug().Str("provider", "openai").Str("model", o.Model).Int("diffBytes", len(diff)).Msg("running openai commit message request")

	started := time.Now()
	client := openai.NewClient()

	model := o.Model
	if model == "" {
		model = defaultOpenAIModel
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:       model,
		Temperature: openai.Float(0.2),
	})
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("openai request timed out after %s: %w", time.Since(started).Truncate(time.Second), ctx.Err())
		}
		return "", fmt.Errorf("failed to generate content from openai: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("empty response from openai")
	}

	result := strings.TrimSpace(chatCompletion.Choices[0].Message.Content)
	if result == "" {
		return "", fmt.Errorf("empty response from openai")
	}

	logger.Debug().Str("provider", "openai").Dur("elapsed", time.Since(started)).Int("messageBytes", len(result)).Msg("finished openai commit message request")

	return result, nil
}
