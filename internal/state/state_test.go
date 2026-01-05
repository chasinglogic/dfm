package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDfmDirCreatesDirectory(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	dir, err := DfmDir()
	if err != nil {
		t.Fatalf("DfmDir returned error: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat on DfmDir path failed: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("DfmDir = %q is not a directory", dir)
	}
}

func TestProfilesAndModulesDir(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	profilesDir, err := ProfilesDir()
	if err != nil {
		t.Fatalf("ProfilesDir returned error: %v", err)
	}

	modulesDir, err := ModulesDir()
	if err != nil {
		t.Fatalf("ModulesDir returned error: %v", err)
	}

	for _, dir := range []string{profilesDir, modulesDir} {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("stat on %q failed: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", dir)
		}
	}
}

func TestLoadCreatesEmptyStateWhenNoFile(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	State = nil

	dir, err := DfmDir()
	if err != nil {
		t.Fatalf("DfmDir returned error: %v", err)
	}

	statePath := filepath.Join(dir, "state.json")
	_ = os.Remove(statePath)

	if err := Load(); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if State == nil {
		t.Fatal("State was nil after Load")
	}

	if State.CurrentProfile != "" {
		t.Fatalf("CurrentProfile = %q, want empty", State.CurrentProfile)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	State = &appState{CurrentProfile: "/tmp/profile"}
	if err := Save(); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	State = nil
	if err := Load(); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if State.CurrentProfile != "/tmp/profile" {
		t.Fatalf("CurrentProfile after Load = %q, want %q", State.CurrentProfile, "/tmp/profile")
	}
}