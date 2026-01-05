package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanDeadSymlinksRemovesOnlyManagedLinks(t *testing.T) {
	root := t.TempDir()
	dfmRoot := t.TempDir()

	managedTarget := filepath.Join(dfmRoot, "profiles", "foo")
	unmanagedTarget := filepath.Join(t.TempDir(), "other", "bar")

	managedLink := filepath.Join(root, "managed")
	unmanagedLink := filepath.Join(root, "unmanaged")

	if err := os.Symlink(managedTarget, managedLink); err != nil {
		t.Fatalf("failed to create managed symlink: %v", err)
	}

	if err := os.Symlink(unmanagedTarget, unmanagedLink); err != nil {
		t.Fatalf("failed to create unmanaged symlink: %v", err)
	}

	if err := cleanDeadSymlinks(root, dfmRoot); err != nil {
		t.Fatalf("cleanDeadSymlinks returned error: %v", err)
	}

	if _, err := os.Lstat(managedLink); !os.IsNotExist(err) {
		t.Fatalf("expected managed dead link to be removed, got err=%v", err)
	}

	if _, err := os.Lstat(unmanagedLink); err != nil {
		t.Fatalf("expected unmanaged link to remain, got err=%v", err)
	}
}