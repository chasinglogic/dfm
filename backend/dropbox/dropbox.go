package dropbox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/dfm/config"
	"gopkg.in/urfave/cli.v1"
)

// Backend implements backend.Backend for a dropbox based remote.
type Backend struct{}

func getDropboxDir() string {
	etc, ok := config.CONFIG.Etc["DROPBOX_DIR"]

	defaultDir := filepath.Join(os.Getenv("HOME"), "Dropbox")
	_, err := os.Stat(defaultDir)
	if err != nil && os.IsNotExist(err) && !ok {
		fmt.Println("Default Dropbox location found.")
		fmt.Println("Set DROPBOX_DIR in your config's etc section.")
		fmt.Println(`Example: 
{
    "Etc": {
         "DROPBOX_DIR": "<path to dropbox folder>
    }
}
`)
		os.Exit(1)
	}

	var ed *string

	if ed, ok = etc.(*string); !ok {
		fmt.Println("Error: Etc DROPBOX_DIR is not correct type.")
		fmt.Println("Got:", etc)
		os.Exit(1)
	}

	defaultDir = *ed
	return defaultDir
}

// Init determines where the Dropbox folder is and sets up dfm
func (b Backend) Init() error {
	config.CONFIG.ConfigDir = filepath.Join(getDropboxDir(), "dfm")
	return nil
}

// Sync has nothing to do. Dropbox handles it all
func (b Backend) Sync(userDir string) error { return nil }

// NewProfile has nothing to do. Dropbox handles it all
func (b Backend) NewProfile(userDir string) error { return nil }

// Commands has nothing to do. Dropbox handles it all
func (b Backend) Commands() []cli.Command { return nil }
