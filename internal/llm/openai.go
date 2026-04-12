package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
)

const defaultOpenAIModel = "gpt-4.1-mini"

// OpenAIProvider generates commit messages using the OpenAI API.
// Requires the OPENAI_API_KEY environment variable to be set.
type OpenAIProvider struct {
	Model string
}

func (o *OpenAIProvider) GenerateCommitMessage(diff string, promptTemplate string) (string, error) {
	client := openai.NewClient()

	model := o.Model
	if model == "" {
		model = defaultOpenAIModel
	}

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

	chatCompletion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:       model,
		Temperature: openai.Float(0.2),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate content from openai: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("empty response from openai")
	}

	result := strings.TrimSpace(chatCompletion.Choices[0].Message.Content)
	if result == "" {
		return "", fmt.Errorf("empty response from openai")
	}

	return result, nil
}
