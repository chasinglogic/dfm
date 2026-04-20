package llm

import (
	"context"
	"strings"
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
func GenerateCommitMessage(ctx context.Context, diff string, providerName string, model string, promptTemplate string) (string, error) {
	provider, err := NewProvider(providerName, model)
	if err != nil {
		return "", err
	}

	return provider.GenerateCommitMessage(ctx, diff, promptTemplate)
}

func buildCommitMessagePrompt(diff string, promptTemplate string) string {
	if strings.TrimSpace(promptTemplate) == "" {
		promptTemplate = defaultCommitMessagePrompt
	}

	return strings.TrimSpace(promptTemplate) + "\n\nDiff:\n" + diff
}
