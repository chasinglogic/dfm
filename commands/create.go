package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// Create will clone the given git repo to the profiles directory, it optionally
// will call link or use depending on the flag given.
func Create(c *cli.Context) error {
	setGlobalOptions(c.Parent())

	var aliasDir string

	if alias := c.String("alias"); alias != "" {
		aliasDir = filepath.Join(c.Parent().String("config"), "profiles", alias)
	}

	url, user := createURL(strings.Split(c.Args().First(), "/"))
	userDir := filepath.Join(c.Parent().String("config"), "profiles", user)
	if cloneErr := cloneRepo(url, user, userDir); cloneErr != nil {
		return cloneErr
	}

	// Just create a symlink in configDir/profiles/ to the other profile name
	if aliasDir != "" {
		if err := os.Symlink(userDir, aliasDir); err != nil {
			fmt.Println("Error creating alias", err, "skipping...")
		}
	}

	if c.Bool("link") {
		return Link(c)
	}

	if c.Bool("use") {
		return Use(c)
	}

	return nil
}

func createURL(s []string) (string, string) {
	if len(s) == 3 {
		return fmt.Sprintf("https://%s", strings.Join(s, "/")), s[1]
	}

	return fmt.Sprintf("https://github.com/%s", strings.Join(s, "/")), s[0]
}

func cloneRepo(url, user, userDir string) error {
	if VERBOSE {
		fmt.Printf("Creating profile in %s\n", userDir)
	}

	c := exec.Command("git", "clone", url, userDir)
	_, err := c.CombinedOutput()
	if err != nil && err.Error() == "exit status 128" {
		return cli.NewExitError("Profile exists, perhaps you meant dfm update?", 128)
	}

	return err
}
