// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

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
