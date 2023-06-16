package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/chasinglogic/dfm/profiles"
	"github.com/spf13/cobra"
)

var (
	link = false
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a git repo as a dotfile profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cloneUrl := args[0]
		if !strings.HasPrefix(cloneUrl, "http") && !strings.HasPrefix(cloneUrl, "git") {
			if !strings.Contains(cloneUrl, "/") {
				cloneUrl = fmt.Sprintf("%s/dotfiles", cloneUrl)
			}

			cloneUrl = fmt.Sprintf("https://github.com/%s", cloneUrl)
		}

		var name string

		if strings.HasPrefix(cloneUrl, "http") {
			split := strings.Split(cloneUrl, "/")
			name = split[len(split)-2]
		} else {
			split := strings.Split(cloneUrl, ":")
			split = strings.Split(split[1], "/")
			name = split[len(split)-2]
		}

		targetPath := path.Join(profiles.ProfileDir, name)

		fmt.Println("cloning profile", cloneUrl, "as", name)
		git := exec.Command("git", "clone", cloneUrl, targetPath)
		git.Stdin = os.Stdin
		git.Stdout = os.Stdout
		git.Stderr = os.Stderr
		git.Start()
		err := git.Wait()
		if err != nil {
			return err
		}

		if !link {
			return nil
		}

		profile, err := loadProfileByName(name)
		if err != nil {
			return err
		}

		err = profile.Link(profiles.LinkOptions{
			Overwrite: linkOverwrite,
		})
		if err != nil {
			return err
		}

		state.CurrentProfile = profile.Name()
		return nil
	},
}

func init() {
	cloneCmd.Flags().BoolVarP(&link, "link", "l", false, "immediately link the profile after cloning")
	cloneCmd.Flags().BoolVarP(&linkOverwrite, "overwrite", "o", false, "remove regular files if they exist at the link path")
}
