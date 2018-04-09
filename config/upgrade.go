package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/dotfiles"
)

type oldConfigJSON struct {
	Backend        string
	ConfigDir      string
	CurrentProfile string
	Etc            map[string]interface{}
}

func checkForOldConfig() bool {
	cfgFile := filepath.Join(os.Getenv("HOME"), ".dfm")
	_, err := os.Stat(cfgFile)
	return err == nil
}

func upgradeConfig() {
	fmt.Println("It looks like you have an old style DFM config.")
	fmt.Println("One second while we upgrade it for you...")

	var oldCfg oldConfigJSON

	cfgFile := filepath.Join(os.Getenv("HOME"), ".dfm")
	jsonBytes, _ := ioutil.ReadFile(cfgFile)

	err := json.Unmarshal(jsonBytes, &oldCfg)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}

	profiles, err := ioutil.ReadDir(
		filepath.Join(oldCfg.ConfigDir, "profiles"))
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	global = Config{
		Dir:                oldCfg.ConfigDir,
		DefaultBackend:     oldCfg.Backend,
		CurrentProfileName: oldCfg.CurrentProfile,
		Profiles:           []dotfiles.Profile{},
	}

	for _, profile := range profiles {
		if strings.HasPrefix(profile.Name(), ".") {
			continue
		}

		global.Profiles = append(
			global.Profiles,
			dotfiles.Profile{
				Name:    profile.Name(),
				Backend: global.DefaultBackend,
				Locations: []string{
					filepath.Join(oldCfg.ConfigDir, "profiles", profile.Name()),
				},
			},
		)
	}

	ans := getInput("Upgrade complete. Should we remove the old config? ")
	if strings.HasPrefix(strings.ToLower(ans), "y") {
		fmt.Println("Removing...")
		err = os.Remove(cfgFile)
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	}

	fmt.Println("Saving new config...")
	err = SaveConfig()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
