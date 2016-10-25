package dfm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

// Config holds some basic information about the state of dfm.
type Config struct {
	Verbose        bool
	ConfigDir      string
	CurrentProfile string
}

// LoadConfig runs as a hook when dfm runs to load the config object from
// ~/.dfm
func LoadConfig(c *cli.Context) error {
	configJSON, rerr := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".dfm"))
	if rerr == os.ErrNotExist {
		configJSON = []byte("{\"CurrentProfile\":false}")
	} else if rerr != nil {
		return rerr
	}

	err := json.Unmarshal(configJSON, &CONFIG)
	if err != nil {
		return err
	}

	if (CONFIG.ConfigDir == "" || CONFIG.ConfigDir == DefaultConfigDir()) &&
		c.String("config") != DefaultConfigDir() {
		CONFIG.ConfigDir = c.String("config")
	}

	if !CONFIG.Verbose && c.Bool("verbose") {
		CONFIG.Verbose = c.Bool("verbose")
	}

	DRYRUN = c.Bool("dry-run")

	return nil
}

// SaveConfig will save the config to the configDir/config.json
func SaveConfig(c *cli.Context) error {
	CONFIG.Verbose = false
	JSON, merr := json.MarshalIndent(CONFIG, "", "\t")
	if merr != nil {
		return merr
	}

	return ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".dfm"), JSON, 0644)
}

func getProfileDir() string {
	return filepath.Join(CONFIG.ConfigDir, "profiles")
}

// DefaultConfigDir will return the default location to store profiles which is
// $XDG_CONFIG_HOME/dfm or $HOME/.config/dfm
func DefaultConfigDir() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdg, "dfm")
}
