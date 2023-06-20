package profiles

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/chasinglogic/dfm/logger"
)

var doSkip = true

var DEFAULT_MAPPINGS = []Mapping{
	{
		Match: "\\/README(\\.md|\\.txt|\\.rst|\\.org)?$",
		Skip:  &doSkip,
	},
	{
		Match: "\\/\\.git\\/",
		Skip:  &doSkip,
	},
	{
		Match: "\\/\\.gitignore$",
		Skip:  &doSkip,
	},
	{
		Match: "\\/LICENSE(\\.md)?$",
		Skip:  &doSkip,
	},
	{
		Match: "\\/\\.dfm\\.yml$",
		Skip:  &doSkip,
	},
	// TODO: Support destination mappings
	// {
	//   Match: "\/\.ggitignore$",
	//   Destination: ".gitignore",
	// },
}

type Mapping struct {
	Match     string   `yaml:"match"`
	TargetDir *string  `yaml:"target_dir"`
	Skip      *bool    `yaml:"skip"`
	TargetOs  []string `yaml:"target_os"`
	LinkAsDir *bool    `yaml:"link_as_dir"`

	rgx *regexp.Regexp
}

func (m Mapping) Matches(path string) bool {
	return m.getRegex().MatchString(path)
}

func (m Mapping) getRegex() *regexp.Regexp {
	if m.rgx == nil {
		var err error
		m.rgx, err = regexp.Compile(m.Match)
		if err != nil {
			logger.Warn.Printf("cannot compile mapping regex: '%s' reason: %s\n", m.Match, err)
			m.rgx, _ = regexp.Compile("^$")
		}

	}

	return m.rgx
}

func (m Mapping) AppliesForCurrentOS() bool {
	if m.TargetOs == nil {
		return true
	}

	for _, osName := range m.TargetOs {
		if runtime.GOOS == strings.ToLower(osName) {
			return true
		}
	}

	return false
}

func (m Mapping) ShouldSkip() bool {
	switch {
	case !m.AppliesForCurrentOS() && m.Skip == nil:
		return true
	case !m.AppliesForCurrentOS():
		return false
	case m.Skip == nil:
		return false
	default:
		return *m.Skip
	}
}

func (m Mapping) ShouldLinkAsDir() bool {
	switch {
	case !m.AppliesForCurrentOS():
		return false
	case m.LinkAsDir == nil:
		return false
	default:
		return *m.LinkAsDir
	}
}
