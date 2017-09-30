package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/config"
	"github.com/spf13/cobra"
)

var (
	overwrite bool
	link      bool
	alias     string
)

func init() {
	Clone.Flags().BoolVarP(&overwrite, "overwrite", "o", false,
		"whether dfm should remove files that exist where a link should go if --link is given")
	Clone.Flags().BoolVarP(&link, "link", "l", false,
		"whether the profile should be linked after being cloned")
	Clone.Flags().StringVarP(&alias, "alias", "a", "",
		"whether the profile should be aliased to a different name")
}

// Clone will clone the given git repo to the profiles directory, it optionally
// will call link or use depending on the flag given.
var Clone = &cobra.Command{
	Use:   "clone",
	Short: "git clone an existing profile from `URL`",
	Run: func(cmd *cobra.Command, args []string) {
		url, user := CreateURL(strings.Split(args[0], "/"))
		userDir := filepath.Join(config.ProfileDir(), user)
		if err := CloneRepo(url, userDir); err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		// Just create a symlink in configDir/profiles/ to the other profile name
		if alias != "" {
			aliasDir := filepath.Join(config.ProfileDir(), alias)
			if err := os.Symlink(userDir, aliasDir); err != nil {
				fmt.Println("Error creating alias", err, "skipping...")
			}
		}

		if link {
			args := []string{"dfm", "link", user}
			if overwrite {
				args = []string{"dfm", "link", "-o", user}
			}

			c := exec.Command(args[0], args[1:]...)
			_, err := c.CombinedOutput()
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}
		}
	},
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
	fmt.Printf("Creating profile in %s\n", profileDir)

	c := exec.Command("git", "clone", url, profileDir)
	_, err := c.CombinedOutput()
	if err != nil && err.Error() == "exit status 128" {
		fmt.Println("Profile exists, perhaps you meant dfm update or link?")
		os.Exit(128)
	}

	return err
}
