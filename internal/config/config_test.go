package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chasinglogic/dfm/internal/state"
)

func TestRepoToNameStripsGitExtension(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"dotfiles":                                "dotfiles",
		"dotfiles.git":                            "dotfiles",
		"git@github.com:chasinglogic/dfm.git":     "dfm",
		"https://github.com/chasinglogic/dfm.git": "dfm",
	}

	for input, want := range cases {
		got := RepoToName(input)
		if got != want {
			t.Fatalf("RepoToName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGetDotfileDirectory(t *testing.T) {
	t.Parallel()

	cfg := &Config{Location: "/tmp/profile"}
	if got, want := cfg.GetDotfileDirectory(), "/tmp/profile"; got != want {
		t.Fatalf("GetDotfileDirectory() = %q, want %q", got, want)
	}

	cfg.RootDir = "dotfiles"
	if got, want := cfg.GetDotfileDirectory(), filepath.Join("/tmp/profile", "dotfiles"); got != want {
		t.Fatalf("GetDotfileDirectory() with RootDir = %q, want %q", got, want)
	}
}

func TestLoadReturnsDefaultWhenFileDoesNotExist(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configFile := filepath.Join(dir, ".dfm.yml")

	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load returned error for missing file: %v", err)
	}

	if cfg.Location != dir {
		t.Fatalf("Location = %q, want %q", cfg.Location, dir)
	}
}

func TestLoadNormalizesModuleLocations(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	dir := t.TempDir()
	configFile := filepath.Join(dir, ".dfm.yml")

	content := []byte(`modules:
  - repository: https://example.com/foo.git
`)
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if len(cfg.Modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(cfg.Modules))
	}

	modulesDir, err := state.ModulesDir()
	if err != nil {
		t.Fatalf("ModulesDir returned error: %v", err)
	}

	expected := filepath.Join(modulesDir, "foo")
	if cfg.Modules[0].Location != expected {
		t.Fatalf("module Location = %q, want %q", cfg.Modules[0].Location, expected)
	}
}