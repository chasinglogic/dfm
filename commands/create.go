package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/common"
	"github.com/urfave/cli"
)

func Create(c *cli.Context) error {
	url, user := createURL(strings.Split(c.Args().First(), "/"))
	userDir := filepath.Join(c.Parent().String("config"), "profiles", user)
	if cloneErr := cloneRepo(url, user, userDir); cloneErr != nil {
		return cloneErr
	}

	links := common.GenerateSymlinks(userDir)
	return common.CreateSymlinks(links)
}

func createURL(s []string) (string, string) {
	if len(s) == 3 {
		return fmt.Sprintf("https://%s", strings.Join(s, "/")), s[1]
	}

	return fmt.Sprintf("https://github.com/%s", strings.Join(s, "/")), s[0]
}

func cloneRepo(url, user, userDir string) error {
	fmt.Printf("Creating profile in %s\n", userDir)
	c := exec.Command("git", "clone", url, userDir)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}
