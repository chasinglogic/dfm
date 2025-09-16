package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/internal/mapping"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/goccy/go-yaml"
)

type LinkMode string

type Config struct {
	Location string `yaml:"-"`

	Repo                   string             `yaml:"repository"`
	RootDir                string             `yaml:"root_dir"`
	PromptForCommitMessage bool               `yaml:"prompt_for_commit_message"`
	Modules                []Config           `yaml:"modules"`
	Mappings               []*mapping.Mapping `yaml:"mappings"`
	LinkMode               string             `yaml:"link_mode"`
}

func Load(configFile string) (*Config, error) {
	config := Config{
		Location: filepath.Dir(configFile),
	}

	content, err := os.ReadFile(configFile)
	if os.IsNotExist(err) {
		return &config, nil
	} else if err != nil {
		return &config, err
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		return &config, err
	}

	modulesDir, err := state.ModulesDir()
	if err != nil {
		return &config, err
	}

	for idx := range config.Modules {
		if config.Modules[idx].Location != "" {
			continue
		}

		config.Modules[idx].Location = filepath.Join(
			modulesDir,
			RepoToName(config.Modules[idx].Repo),
		)
	}

	return &config, nil
}

func RepoToName(repo string) string {
	return strings.ReplaceAll(filepath.Base(repo), ".git", "")
}
