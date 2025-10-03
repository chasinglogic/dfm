package mapping

import (
	"encoding/json"
	"regexp"
	"runtime"
	"strings"
)

type Action int

func (a Action) String() string {
	switch a {
	case ActionNone:
		return "NONE"
	case ActionSkip:
		return "SKIP"
	case ActionTranslate:
		return "TRANSLATE"
	case ActionLinkAsDir:
		return "LINK_AS_DIR"
	default:
		return "UNKNOWN"
	}
}

const (
	ActionNone = iota
	ActionSkip
	ActionTranslate
	ActionLinkAsDir
)

type Mapping struct {
	rgx *regexp.Regexp `yaml:"-" json:"-"`

	Match     string `yaml:"match"`
	LinkAsDir bool   `yaml:"link_as_dir"`
	Skip      bool   `yaml:"skip"`
	Dest      string `yaml:"dest"`
	TargetOS  string `yaml:"target_os"`
}

func (m *Mapping) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func (m *Mapping) IsMatch(path string) bool {
	if m.rgx == nil {
		m.rgx = regexp.MustCompile(m.Match)
	}

	matchesRegex := m.rgx.Match([]byte(path))
	isTargetOS := strings.ToLower(m.TargetOS) == runtime.GOOS || m.TargetOS == ""
	return matchesRegex && isTargetOS
}

func (m *Mapping) Action() Action {
	if m.Skip {
		return ActionSkip
	}

	if m.LinkAsDir {
		return ActionLinkAsDir
	}

	if m.Dest != "" {
		return ActionTranslate
	}

	return ActionNone
}
