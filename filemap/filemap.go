package filemap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Mappings is a list of Mapping's
type Mappings []Mapping

// Mapping is a way for DFM to map custom locations to custom locations
type Mapping struct {
	Match  string `yaml:"match"`
	Dest   string `yaml:"dest"`
	Regexp bool   `yaml:"regexp"`
	IsDir  bool   `yaml:"isDir"`
	Skip   bool   `yaml:"skip"`
}

func New() Mapping {
	return Mapping{
		Match:  "",
		Dest:   "",
		Regexp: false,
		IsDir:  true,
		Skip:   false,
	}
}

func (m Mapping) Matches(filename string) bool {
	if m.Regexp {
		rg, err := regexp.Compile(m.Match)
		if err != nil {
			fmt.Println("ERROR compiling match regex:", err.Error())
			return false
		}

		return rg.Match([]byte(filename))
	}

	return strings.HasPrefix(filename, m.Match)
}

func (m Mappings) Matches(filename string) *Mapping {
	for i := range m {
		if m[i].Matches(filename) {
			return &m[i]
		}
	}

	return nil
}

func DefaultMappings() Mappings {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".config")
	}

	xdgconfig := New()
	xdgconfig.Match = "^[.]?config"
	xdgconfig.Regexp = true
	xdgconfig.Dest = xdg

	return Mappings{
		xdgconfig,
		{
			Match:  "^[.]?ggitignore",
			IsDir:  false,
			Regexp: true,
			Dest:   "gitignore",
		},
		{
			Match:  "^\\.git",
			IsDir:  true,
			Regexp: true,
			Skip:   true,
		},
		{
			Match:  "^\\.gitignore$",
			IsDir:  false,
			Regexp: true,
			Skip:   true,
		},
		{
			Match:  "^LICENSE(\\.md)?$",
			IsDir:  false,
			Regexp: true,
			Skip:   true,
		},
		{
			Match:  "^\\.dfm\\.yml$",
			IsDir:  false,
			Regexp: true,
			Skip:   true,
		},
		{
			Match:  "^README(\\.md)?$",
			IsDir:  false,
			Regexp: true,
			Skip:   true,
		},
	}
}
