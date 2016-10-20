package dfm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

// Create will clone the given git repo to the profiles directory, it optionally
// will call link or use depending on the flag given.
func Create(c *cli.Context) error {
	_, cerr := loadConfig(c.Parent())
	if cerr != nil {
		return cli.NewExitError(cerr.Error(), 3)
	}

	var aliasDir string

	if alias := c.String("alias"); alias != "" {
		aliasDir = filepath.Join(getProfileDir(c), alias)
	}

	url, user := CreateURL(strings.Split(c.Args().First(), "/"))
	userDir := filepath.Join(getProfileDir(c), user)
	if cloneErr := CloneRepo(url, user, userDir); cloneErr != nil {
		return cloneErr
	}

	// Just create a symlink in configDir/profiles/ to the other profile name
	if aliasDir != "" && aliasDir != nil {
		if err := os.Symlink(userDir, aliasDir); err != nil {
			fmt.Println("Error creating alias", err, "skipping...")
		}
	}

	if c.Bool("link") {
		return Link(c)
	}

	return nil
}

func CreateURL(s []string) (string, string) {
	if len(s) == 2 {
		return fmt.Sprintf("https://github.com/%s", strings.Join(s, "/")), s[0]
	}

	return strings.Join(s, "/"), s[len(s)-2]
}

func CloneRepo(url, user, userDir string) error {
	if VERBOSE {
		fmt.Printf("Creating profile in %s\n", userDir)
	}

	c := exec.Command("git", "clone", url, userDir)
	_, err := c.CombinedOutput()
	if err != nil && err.Error() == "exit status 128" {
		return cli.NewExitError("Profile exists, perhaps you meant dfm update or link?", 128)
	}

	return err
}
