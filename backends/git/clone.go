package dfm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

// Clone will clone the given git repo to the profiles directory, it optionally
// will call link or use depending on the flag given.
func Clone(c *cli.Context) error {
	var aliasDir string

	if alias := c.String("alias"); alias != "" {
		aliasDir = filepath.Join(getProfileDir(), alias)
	}

	url, user := CreateURL(strings.Split(c.Args().First(), "/"))
	userDir := filepath.Join(getProfileDir(), user)
	if err := CloneRepo(url, userDir); err != nil {
		return err
	}

	// Just create a symlink in configDir/profiles/ to the other profile name
	if aliasDir != "" {
		if err := os.Symlink(userDir, aliasDir); err != nil {
			fmt.Println("Error creating alias", err, "skipping...")
		}
	}

	if c.Bool("link") {
		err := CreateSymlinks(userDir, os.Getenv("HOME"), c.Bool("overwrite"))
		if err != nil {
			return err
		}

		CONFIG.CurrentProfile = user
	}

	return nil
}

// CreateURL will add the missing github.com for the shorthand version of
// links.
func CreateURL(s []string) (string, string) {
	if len(s) == 2 {
		return fmt.Sprintf("https://github.com/%s", strings.Join(s, "/")), s[0]
	}

	return strings.Join(s, "/"), s[len(s)-2]
}

// CloneRepo will git clone the provided url into the appropriate profileDir
func CloneRepo(url, profileDir string) error {
	if CONFIG.Verbose {
		fmt.Printf("Creating profile in %s\n", profileDir)
	}

	c := exec.Command("git", "clone", url, profileDir)
	_, err := c.CombinedOutput()
	if err != nil && err.Error() == "exit status 128" {
		return cli.NewExitError("Profile exists, perhaps you meant dfm update or link?", 128)
	}

	return err
}
