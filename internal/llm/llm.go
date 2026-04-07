package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const defaultCommitMessagePrompt = `You write git commit messages for
configuration-only diffs.

Rules:
- Describe only configuration changes shown in the diff.
- Do not speculate about runtime impact, intent, bug fixes, or user-facing
  behavior unless explicitly shown.
- Prefer concrete config terms such as key, value, flag, path, default,
  threshold, enabled, disabled, renamed, or removed.
- Avoid vague words like "improve", "update", "refactor", or "fix" unless the
  diff clearly supports them.

Output format:
- First line: concise imperative subject, max 72 characters.
- No trailing period.
- Use a body only if the diff includes more than one distinct config change.
- If a body is present, add exactly one blank line, then "- " bullets.
- Each bullet must describe one concrete config change from the diff.
- Return only the raw commit message text. Do not use markdown fences or extra
  commentary.`

// GenerateCommitMessage generates a commit message based on a git diff using the specified provider.
func GenerateCommitMessage(diff string, provider string, promptTemplate string) (string, error) {
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

	prompt := buildCommitMessagePrompt(diff, promptTemplate)

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

func buildCommitMessagePrompt(diff string, promptTemplate string) string {
	if strings.TrimSpace(promptTemplate) == "" {
		promptTemplate = defaultCommitMessagePrompt
	}

	return strings.TrimSpace(promptTemplate) + "\n\nDiff:\n" + diff
}
