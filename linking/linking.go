package linking

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/config"
	"github.com/chasinglogic/dfm/filemap"
)

type Linker interface {
	Link(mappings filemap.Mappings, config Config) error
}

// Config changes the behavior of linking functions, it acts as a
// store for common functoin arguments
type Config struct {
	DryRun    bool
	Overwrite bool
}

// LinkInfo holds the src and destination for our symlink.
type LinkInfo struct {
	Src  string
	Dest string
}

func (l *LinkInfo) String() string {
	return fmt.Sprintf("%s -> %s", l.Dest, l.Src)
}

func (l *LinkInfo) Link(config Config) error {
	if config.DryRun {
		fmt.Println("Creating symlink", l)
	}

	e := removeIfNeeded(l, config)
	if e != nil && strings.Contains(e.Error(), "cowardly") {
		fmt.Println(e)
	} else if e != nil {
		return e
	}

	if !config.DryRun {
		if e := os.Symlink(l.Src, l.Dest); e != nil {
			return e
		}
	}

	return nil
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
func GenerateSymlink(sourceDir, targetDir string, file os.FileInfo) LinkInfo {
	target := getTargetName(file.Name())

	if targetDir != os.Getenv("HOME") {
		target = file.Name()
	}

	ln := LinkInfo{
		filepath.Join(sourceDir, file.Name()),
		filepath.Join(targetDir, target),
	}

	return ln
}

// removeIfNeeded will check if the link destination exists and delete it if
// appropriate.
func removeIfNeeded(link *LinkInfo, config Config) error {
	info, err := os.Lstat(link.Dest)
	if err == nil && (config.Overwrite || info.Mode()&os.ModeSymlink == os.ModeSymlink) {
		if config.DryRun {
			fmt.Printf("%s already exists, removing.\n", link.Dest)
		}

		if !config.DryRun {
			if rmerr := os.Remove(link.Dest); rmerr != nil {
				return fmt.Errorf("Unable to remove %s: %s",
					link.Dest,
					rmerr.Error())
			}
		}

	} else if err == nil {
		msg := fmt.Sprintf(
			"%s already exists and is not a symlink, cowardly refusing to remove",
			link.Dest,
		)
		return errors.New(msg)
	}

	return nil
}

// CreateSymlinks will read all of the files at sourceDir and link them to the
// appropriate location in targetDir, if there is a folder named config in
// sourceDir CreateSymlinks will run itself using that folder as sourceDir and
// targetDir as XDG_config.CONFIG_HOME or HOME/.config if XDG_config.CONFIG_HOME is not set.
func CreateSymlinks(sourceDir, targetDir string, config Config, mappings filemap.Mappings) error {
	sourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return err
	}

	links, err := GenerateSymlinks(sourceDir, targetDir, mappings)
	if err != nil {
		return err
	}

	for _, link := range links {
		err := link.Link(config)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateSymlinks will create the symlinks so we know what they were supposed
// to be prior to removing the profile.
func GenerateSymlinks(profileDir, target string, mappings filemap.Mappings) ([]LinkInfo, error) {
	var lnks []LinkInfo

	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		return lnks, err
	}

	for _, file := range files {
		// Handle XDG_config.CONFIG_HOME special case.
		mapping := mappings.Matches(file.Name())
		if mapping != nil {
			if mapping.Skip {
				continue
			} else if mapping.IsDir {
				newTargetDir := config.ExpandFilePath(mapping.Dest)
				multiLinks, err := GenerateSymlinks(filepath.Join(profileDir, file.Name()), newTargetDir, mappings)
				if err != nil {
					return lnks, err
				}

				lnks = append(lnks, multiLinks...)
			} else {
				link := GenerateSymlink(profileDir, mapping.Dest, file)
				lnks = append(lnks, link)
			}

			continue
		}

		link := GenerateSymlink(profileDir, target, file)
		lnks = append(lnks, link)
	}

	return lnks, nil
}
