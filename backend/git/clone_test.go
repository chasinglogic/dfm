package git

import (
	"strings"
	"testing"
)

func TestCreateLongURL(t *testing.T) {
	url, user := CreateURL(strings.Split("https://bitbucket.org/chasinglogic/dotfiles", "/"))
	if url != "https://bitbucket.org/chasinglogic/dotfiles" {
		t.Errorf("Expected: %s Got: %s", "https://bitbucket.org/chasinglogic/dotfiles", url)
	}

	if user != "chasinglogic" {
		t.Errorf("Expected: %s Got: %s", "chasinglogic", user)
	}
}

func TestCreateShortURL(t *testing.T) {
	url, user := CreateURL(strings.Split("chasinglogic/dotfiles", "/"))
	if url != "https://github.com/chasinglogic/dotfiles" {
		t.Errorf("Expected: %s Got: %s", "https://github.com/chasinglogic/dotfiles", url)
	}

	if user != "chasinglogic" {
		t.Errorf("Expected: %s Got: %s", "chasinglogic", user)
	}
}

func TestCreateSSHURL(t *testing.T) {
	url, user := CreateURL(strings.Split("git@github.com:/chasinglogic/dotfiles", "/"))
	if url != "git@github.com:/chasinglogic/dotfiles" {
		t.Errorf("Expected: %s Got: %s", "git@github.com:/chasinglogic/dotfiles", url)
	}

	if user != "chasinglogic" {
		t.Errorf("Expected: %s Got: %s", "chasinglogic", user)
	}
}
