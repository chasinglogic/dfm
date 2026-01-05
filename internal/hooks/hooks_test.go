package hooks

import "testing"

func TestParseStringHook(t *testing.T) {
	args, err := parse("echo hello")
	if err != nil {
		t.Fatalf("parse returned error: %v", err)
	}

	if len(args) != 3 || args[0] != "/bin/sh" || args[1] != "-c" || args[2] != "echo hello" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestParseMapHook(t *testing.T) {
	hook := map[string]any{
		"interpreter": "bash -c",
		"script":      "echo hello",
	}

	args, err := parse(hook)
	if err != nil {
		t.Fatalf("parse returned error: %v", err)
	}

	if len(args) != 3 || args[0] != "bash" || args[1] != "-c" || args[2] != "echo hello" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestParseInvalidInterpreter(t *testing.T) {
	hook := map[string]any{
		"interpreter": 123,
		"script":      "echo",
	}

	if _, err := parse(hook); err != ErrInvalidHook {
		t.Fatalf("expected ErrInvalidHook, got %v", err)
	}
}

func TestParseInvalidScript(t *testing.T) {
	hook := map[string]any{
		"interpreter": "bash -c",
		"script":      123,
	}

	if _, err := parse(hook); err != ErrInvalidHook {
		t.Fatalf("expected ErrInvalidHook, got %v", err)
	}
}

func TestParseUnsupportedType(t *testing.T) {
	if _, err := parse(123); err != ErrInvalidHook {
		t.Fatalf("expected ErrInvalidHook, got %v", err)
	}
}

func TestHooksExecuteNoHooksDefined(t *testing.T) {
	var h Hooks
	if err := h.Execute(".", "nonexistent"); err != nil {
		t.Fatalf("Execute returned error for missing hook: %v", err)
	}
}

func TestHooksExecuteRunsValidHook(t *testing.T) {
	h := Hooks{
		"test": []any{"true"},
	}

	if err := h.Execute(".", "test"); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}