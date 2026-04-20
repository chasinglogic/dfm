package llm

import (
	"context"
	"fmt"
)

// Provider generates commit messages from diffs using an LLM.
type Provider interface {
	// GenerateCommitMessage takes a diff and a prompt template and returns a commit message.
	GenerateCommitMessage(ctx context.Context, diff string, promptTemplate string) (string, error)
}

// NewProvider returns a Provider for the given provider name and model.
// If model is empty, the provider will use its default model.
//
// Available providers:
//   - "gemini":     Google Gemini API (requires GEMINI_API_KEY)
//   - "gemini-cli": Gemini CLI tool (requires `gemini` in PATH)
//   - "claude":     Claude CLI (requires `claude` in PATH)
//   - "openai":     OpenAI API (requires OPENAI_API_KEY)
//   - "codex":      OpenAI Codex CLI (requires `codex` in PATH)
func NewProvider(providerName string, model string) (Provider, error) {
	switch providerName {
	case "gemini":
		return &GeminiProvider{Model: model}, nil
	case "gemini-cli":
		return &GeminiCLIProvider{Model: model}, nil
	case "claude":
		return &ClaudeProvider{Model: model}, nil
	case "openai":
		return &OpenAIProvider{Model: model}, nil
	case "codex":
		return &CodexProvider{Model: model}, nil
	default:
		return nil, fmt.Errorf("unsupported model provider: %s (supported: gemini, gemini-cli, claude, openai, codex)", providerName)
	}
}
