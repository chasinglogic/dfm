package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const defaultGeminiModel = "gemini-2.5-flash"

// GeminiProvider generates commit messages using Google's Gemini API.
type GeminiProvider struct {
	Model string
}

func (g *GeminiProvider) GenerateCommitMessage(ctx context.Context, diff string, promptTemplate string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	modelName := g.Model
	if modelName == "" {
		modelName = defaultGeminiModel
	}

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.2)

	prompt := buildCommitMessagePrompt(diff, promptTemplate)
	logger.Debug().Str("provider", "gemini").Str("model", modelName).Int("diffBytes", len(diff)).Msg("running gemini commit message request")

	started := time.Now()

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("gemini request timed out after %s: %w", time.Since(started).Truncate(time.Second), ctx.Err())
		}
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		result := strings.TrimSpace(string(text))
		if result == "" {
			return "", fmt.Errorf("empty response from gemini")
		}

		logger.Debug().Str("provider", "gemini").Dur("elapsed", time.Since(started)).Int("messageBytes", len(result)).Msg("finished gemini commit message request")
		return result, nil
	}

	return "", fmt.Errorf("unexpected response format from gemini")
}
