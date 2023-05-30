package profiles

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/chasinglogic/dfm/logger"
	"gopkg.in/yaml.v3"
)

// The XDG spec for Darwin return Application Support as the config directory
// but we want to maintain backwards compatibility with previous dfm versions as
// well as behave more commonly across unixes.
var DFMDir = path.Join(
	os.Getenv("HOME"),
	".config",
	"dfm",
)
var ProfileDir = path.Join(DFMDir, "profiles")
var ModulesDir = path.Join(DFMDir, "modules")

type ProfileConfig struct {
	Name      string `yaml:"name"`
	Location  string `yaml:"location"`
	TargetDir string `yaml:"target_dir"`

	// Repo is here to support either Repo or Repository for backwards
	// compatibility.
	Repo                   string          `yaml:"repo,omitempty"`
	Repository             string          `yaml:"repository"`
	Branch                 string          `yaml:"branch"`
	Hooks                  Hooks           `yaml:"hooks"`
	LinkMode               string          `yaml:"link_mode"`
	Mappings               []Mapping       `yaml:"mappings"`
	Modules                []ProfileConfig `yaml:"modules"`
	PullOnly               bool            `yaml:"pull_only"`
	PromptForCommitMessage bool            `yaml:"prompt_for_commit_message"`
}

func (cfg *ProfileConfig) String() string {
	yml, _ := yaml.Marshal(cfg)
	return string(yml)
}

func (cfg *ProfileConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cfg.PullOnly = false
	cfg.Modules = []ProfileConfig{}
	cfg.TargetDir = os.Getenv("HOME")
	cfg.LinkMode = "after"

	// Prevent infinite recursion by casting to a temporary type.
	type defaultConfig ProfileConfig
	err := unmarshal((*defaultConfig)(cfg))
	if err != nil {
		return err
	}

	if strings.HasPrefix(cfg.Location, "~") {
		cfg.Location = strings.Replace(cfg.Location, "~", os.Getenv("HOME"), 1)
	}

	if cfg.Location == "" {
		base := strings.Replace(path.Base(cfg.GetRepo()), ".git", "", -1)
		cfg.Name = base
		cfg.Location = path.Join(ModulesDir, base)
	} else {
		cfg.Name = path.Base(cfg.Location)
	}

	return nil
}

func (cfg *ProfileConfig) GetRepo() string {
	if cfg.Repo != "" {
		return cfg.Repo
	}

	return cfg.Repository
}

func DefaultConfig(where string) ProfileConfig {
	return ProfileConfig{
		Name:      path.Base(where),
		Location:  where,
		TargetDir: os.Getenv("HOME"),

		PullOnly: false,
		Modules:  []ProfileConfig{},
		Branch:   "",
	}
}

// Load will load the .dfm.yml file in where and create a Profile object
// representing this location.
func Load(where string) (Profile, error) {
	p := DefaultConfig(where)

	logger.Debug.Printf("Loading profile from: %s\n", where)

	dotdfm := path.Join(where, ".dfm.yml")
	if _, err := os.Stat(dotdfm); err == nil {
		logger.Debug.Println("Profile has a .dfm.yml file.")
		buf, err := os.ReadFile(dotdfm)
		if err != nil {
			return Profile{}, err
		}

		err = yaml.Unmarshal(buf, &p)
		if err != nil {
			return Profile{}, err
		}
	}

	for _, cfg := range p.Modules {
		if cfg.Repository != "" && cfg.Repo != "" {
			panic(
				fmt.Sprintf(
					"both repo and repository where specified, unclear which to use. repo: %s repository: %s",
					cfg.Repo,
					cfg.Repository,
				),
			)
		}

	}

	if p.Location == "" {
		p.Location = where
	}

	return FromConfig(p), nil
}
