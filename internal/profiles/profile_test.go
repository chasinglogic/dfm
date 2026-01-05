package profiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chasinglogic/dfm/internal/config"
	"github.com/chasinglogic/dfm/internal/mapping"
)

func TestLinkCreatesSymlinkInHome(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()

	if err := os.WriteFile(filepath.Join(repo, "foo"), []byte("data"), 0644); err != nil {
		t.Fatalf("failed to write file in repo: %v", err)
	}

	t.Setenv("HOME", home)

	cfg := &config.Config{Location: repo}
	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := p.Link(false); err != nil {
		t.Fatalf("Link returned error: %v", err)
	}

	targetPath := filepath.Join(home, "foo")
	info, err := os.Lstat(targetPath)
	if err != nil {
		t.Fatalf("Lstat on targetPath failed: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("%q is not a symlink", targetPath)
	}

	linkTarget, err := os.Readlink(targetPath)
	if err != nil {
		t.Fatalf("Readlink failed: %v", err)
	}

	expectedTarget := filepath.Join(repo, "foo")
	if linkTarget != expectedTarget {
		t.Fatalf("symlink target = %q, want %q", linkTarget, expectedTarget)
	}
}

func TestLinkSkipsGitAndConfigFiles(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()

	if err := os.Mkdir(filepath.Join(repo, ".git"), 0755); err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repo, ".dfm.yml"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write .dfm.yml: %v", err)
	}

	t.Setenv("HOME", home)

	cfg := &config.Config{Location: repo}
	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := p.Link(false); err != nil {
		t.Fatalf("Link returned error: %v", err)
	}

	if _, err := os.Lstat(filepath.Join(home, ".git")); !os.IsNotExist(err) {
		t.Fatalf("expected no symlink for .git in home, got err=%v", err)
	}

	if _, err := os.Lstat(filepath.Join(home, ".dfm.yml")); !os.IsNotExist(err) {
		t.Fatalf("expected no symlink for .dfm.yml in home, got err=%v", err)
	}
}

func TestLinkAppliesTranslateMapping(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()

	if err := os.WriteFile(filepath.Join(repo, "foo"), []byte("data"), 0644); err != nil {
		t.Fatalf("failed to write file in repo: %v", err)
	}

	customDir := filepath.Join(home, "custom")

	m := &mapping.Mapping{
		Match: "foo$",
		Dest:  customDir,
	}

	t.Setenv("HOME", home)

	cfg := &config.Config{
		Location: repo,
		Mappings: []*mapping.Mapping{m},
	}

	p, err := New(cfg)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err := p.Link(false); err != nil {
		t.Fatalf("Link returned error: %v", err)
	}

	defaultPath := filepath.Join(home, "foo")
	if _, err := os.Lstat(defaultPath); !os.IsNotExist(err) {
		t.Fatalf("expected no default link at %q, got err=%v", defaultPath, err)
	}

	translated := filepath.Join(customDir, "foo")
	info, err := os.Lstat(translated)
	if err != nil {
		t.Fatalf("Lstat on translated path failed: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("%q is not a symlink", translated)
	}

	linkTarget, err := os.Readlink(translated)
	if err != nil {
		t.Fatalf("Readlink failed: %v", err)
	}

	expectedTarget := filepath.Join(repo, "foo")
	if linkTarget != expectedTarget {
		t.Fatalf("symlink target = %q, want %q", linkTarget, expectedTarget)
	}
}

func TestDeleteIfExistsBehavior(t *testing.T) {
	dir := t.TempDir()

	t.Run("no file", func(t *testing.T) {
		if err := deleteIfExists(false, filepath.Join(dir, "does-not-exist")); err != nil {
			t.Fatalf("expected no error when file does not exist, got %v", err)
		}
	})

	t.Run("directory", func(t *testing.T) {
		p := filepath.Join(dir, "subdir")
		if err := os.Mkdir(p, 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}

		if err := deleteIfExists(false, p); err == nil {
			t.Fatalf("expected error when attempting to delete directory")
		}

		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected directory to remain, got err=%v", err)
		}
	})

	t.Run("regular file without overwrite", func(t *testing.T) {
		p := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(p, []byte("data"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		if err := deleteIfExists(false, p); err == nil {
			t.Fatalf("expected error when deleting regular file without overwrite")
		}

		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected file to remain, got err=%v", err)
		}
	})

	t.Run("regular file with overwrite", func(t *testing.T) {
		p := filepath.Join(dir, "file-overwrite.txt")
		if err := os.WriteFile(p, []byte("data"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		if err := deleteIfExists(true, p); err != nil {
			t.Fatalf("expected no error when deleting regular file with overwrite, got %v", err)
		}

		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Fatalf("expected file to be removed, got err=%v", err)
		}
	})
}