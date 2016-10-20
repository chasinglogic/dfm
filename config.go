package dfm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

type Config struct {
	Verbose        bool
	ConfigDir      string
	CurrentProfile string
}

func LoadConfig(c *cli.Context) error {
	configJSON, rerr := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".dfm"))
	if rerr != nil {
		return rerr
	}

	err := json.Unmarshal(configJSON, &CONFIG)
	if err != nil {
		return err
	}

	CONFIG.ConfigDir = c.String("config")
	CONFIG.Verbose = c.Bool("verbose")
	DRYRUN = c.Bool("dry-run")

	return nil
}

// SaveConfig will save the config to the configDir/config.json
func SaveConfig(c *cli.Context) error {
	JSON, merr := json.MarshalIndent(CONFIG, "", "\t")
	if merr != nil {
		return merr
	}

	return ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".dfm"), JSON, 0644)
}

func getProfileDir() string {
	return filepath.Join(CONFIG.ConfigDir, "profiles")
}

func DefaultConfigDir() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdg, "dfm")
}
