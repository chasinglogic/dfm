package llm

import "fmt"

// Provider generates commit messages from diffs using an LLM.
type Provider interface {
	// GenerateCommitMessage takes a diff and a prompt template and returns a commit message.
	GenerateCommitMessage(diff string, promptTemplate string) (string, error)
}

// NewProvider returns a Provider for the given provider name and model.
// If model is empty, the provider will use its default model.
func NewProvider(providerName string, model string) (Provider, error) {
	switch providerName {
	case "gemini":
		return &GeminiProvider{Model: model}, nil
	case "claude":
		return &ClaudeProvider{Model: model}, nil
	case "openai":
		return &OpenAIProvider{Model: model}, nil
	default:
		return nil, fmt.Errorf("unsupported model provider: %s (supported: gemini, claude, openai)", providerName)
	}
}
