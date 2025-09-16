package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/dfm/internal/hooks"
	"github.com/chasinglogic/dfm/internal/mapping"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/goccy/go-yaml"
)

type LinkMode string

type Config struct {
	Location string `yaml:"-"`

	LinkMode               string             `yaml:"link_mode"`
	Mappings               []*mapping.Mapping `yaml:"mappings"`
	Modules                []Config           `yaml:"modules"`
	PromptForCommitMessage bool               `yaml:"prompt_for_commit_message"`
	PullOnly               bool               `yaml:"pull_only"`
	Repo                   string             `yaml:"repository"`
	RootDir                string             `yaml:"root_dir"`
	Hooks                  hooks.Hooks        `yaml:"hooks"`
}

func (c *Config) String() string {
	data, _ := json.Marshal(c)
	return string(data)
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

func (c *Config) GetDotfileDirectory() string {
	// This works because if c.RootDir is "" then filepath.Join ignores it.
	return filepath.Join(c.Location, c.RootDir)
}
