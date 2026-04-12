package llm

import (
	"strings"
	"testing"
)

func TestNewProviderGemini(t *testing.T) {
	p, err := NewProvider("gemini", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*GeminiProvider); !ok {
		t.Fatalf("expected *GeminiProvider, got %T", p)
	}
}

func TestNewProviderGeminiWithModel(t *testing.T) {
	p, err := NewProvider("gemini", "gemini-2.0-flash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gp, ok := p.(*GeminiProvider)
	if !ok {
		t.Fatalf("expected *GeminiProvider, got %T", p)
	}
	if gp.Model != "gemini-2.0-flash" {
		t.Fatalf("expected model gemini-2.0-flash, got %s", gp.Model)
	}
}

func TestNewProviderClaude(t *testing.T) {
	p, err := NewProvider("claude", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*ClaudeProvider); !ok {
		t.Fatalf("expected *ClaudeProvider, got %T", p)
	}
}

func TestNewProviderClaudeWithModel(t *testing.T) {
	p, err := NewProvider("claude", "opus")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, ok := p.(*ClaudeProvider)
	if !ok {
		t.Fatalf("expected *ClaudeProvider, got %T", p)
	}
	if cp.Model != "opus" {
		t.Fatalf("expected model opus, got %s", cp.Model)
	}
}

func TestNewProviderOpenAI(t *testing.T) {
	p, err := NewProvider("openai", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*OpenAIProvider); !ok {
		t.Fatalf("expected *OpenAIProvider, got %T", p)
	}
}

func TestNewProviderOpenAIWithModel(t *testing.T) {
	p, err := NewProvider("openai", "gpt-4o")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	op, ok := p.(*OpenAIProvider)
	if !ok {
		t.Fatalf("expected *OpenAIProvider, got %T", p)
	}
	if op.Model != "gpt-4o" {
		t.Fatalf("expected model gpt-4o, got %s", op.Model)
	}
}

func TestNewProviderUnsupported(t *testing.T) {
	_, err := NewProvider("unsupported", "")
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
	if !strings.Contains(err.Error(), "unsupported model provider") {
		t.Fatalf("unexpected error message: %v", err)
	}
	if !strings.Contains(err.Error(), "gemini, claude, openai") {
		t.Fatalf("error should list supported providers: %v", err)
	}
}
