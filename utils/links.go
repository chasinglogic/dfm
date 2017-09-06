package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// LinkInfo holds the src and destination for our symlink.
type LinkInfo struct {
	Src  string
	Dest string
}

func (l *LinkInfo) String() string {
	return fmt.Sprintf("%s -> %s", l.Dest, l.Src)
}

// getTargetName determines if we need to add a dot to the destination or not.
func getTargetName(n string) string {
	if !strings.HasPrefix(n, ".") {
		return "." + n
	}

	return n
}

// GenerateSymlink will create a LinkInfo with the appropriate destination,
// handling the XDG_config.CONFIG_HOME special case.
func GenerateSymlink(sourceDir, targetDir string, file os.FileInfo, DRYRUN bool) *LinkInfo {
	target := getTargetName(file.Name())

	if strings.HasSuffix(sourceDir, "config") {
		target = file.Name()
	}

	ln := &LinkInfo{
		filepath.Join(sourceDir, file.Name()),
		filepath.Join(targetDir, target),
	}

	if DRYRUN {
		fmt.Printf("Generated symlink %s\n", ln.String())
	}

	return ln
}

// removeIfNeeded will check if the link destination exists and delete it if
// appropriate.
func removeIfNeeded(link *LinkInfo, DRYRUN, overwrite bool) error {
	info, err := os.Lstat(link.Dest)
	if err == nil && (overwrite || info.Mode()&os.ModeSymlink == os.ModeSymlink) {
		if DRYRUN {
			fmt.Printf("%s already exists, removing.\n", link.Dest)
		}

		if !DRYRUN {
			if rmerr := os.Remove(link.Dest); rmerr != nil {
				return fmt.Errorf("Unable to remove %s: %s",
					link.Dest,
					rmerr.Error())
			}
		}

	} else if err == nil {
		return fmt.Errorf("%s already exists and is not a symlink, cowardly refusing to remove", link.Dest)
	}

	return nil
}

// CreateSymlinks will read all of the files at sourceDir and link them to the
// appropriate location in targetDir, if there is a folder named config in
// sourceDir CreateSymlinks will run itself using that folder as sourceDir and
// targetDir as XDG_config.CONFIG_HOME or HOME/.config if XDG_config.CONFIG_HOME is not set.
func CreateSymlinks(sourceDir, targetDir string, DRYRUN, overwrite bool) error {
	sourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range files {
		// Skip the .git directory
		if file.Name() == ".git" {
			continue
		}

		// Handle XDG_config.CONFIG_HOME special case.
		if (file.Name() == "config" || file.Name() == ".config") && file.IsDir() {
			xdg := os.Getenv("XDG_config.CONFIG_HOME")
			if xdg == "" {
				xdg = filepath.Join(os.Getenv("HOME"), ".config")
			}

			err := CreateSymlinks(filepath.Join(sourceDir, file.Name()), xdg,
				DRYRUN, overwrite)
			if err != nil {
				return err
			}

			continue
		}

		link := GenerateSymlink(sourceDir, targetDir, file, DRYRUN)
		e := removeIfNeeded(link, DRYRUN, overwrite)
		if e != nil {
			fmt.Println(e)
			continue
		}

		if DRYRUN {
			fmt.Println("Creating symlink", link)
		}

		if !DRYRUN {
			if err := os.Symlink(link.Src, link.Dest); err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}
