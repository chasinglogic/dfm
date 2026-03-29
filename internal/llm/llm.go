package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateCommitMessage generates a commit message based on a git diff using the specified provider.
func GenerateCommitMessage(diff string, provider string) (string, error) {
	if provider != "gemini" {
		return "", fmt.Errorf("unsupported model provider: %s", provider)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.2)

	prompt := fmt.Sprintf(`You are an expert developer. Generate a concise, conventional git commit message based on the following git diff.
The message should have a short summary line (max 50 characters) followed by a blank line and then a detailed description if necessary.
Return ONLY the raw commit message text without any markdown formatting or extra text. Do not wrap it in backticks.

Diff:
%s`, diff)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		return strings.TrimSpace(string(text)), nil
	}

	return "", fmt.Errorf("unexpected response format from gemini")
}
