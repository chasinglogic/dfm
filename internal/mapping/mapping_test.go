package mapping

import (
	"runtime"
	"testing"
)

func notCurrentOS() string {
	if runtime.GOOS == "linux" {
		return "darwin"
	} else {
		return "linux"
	}
}

func TestIsMatch(t *testing.T) {
	cases := []struct {
		name     string
		mapping  *Mapping
		path     string
		expected bool
	}{
		{
			name: "simple match",
			mapping: &Mapping{
				Match: "foo",
			},
			path:     ".config/bar/foo",
			expected: true,
		},
		{
			name: "simple mismatch",
			mapping: &Mapping{
				Match: "baz",
			},
			path:     ".config/bar/foo",
			expected: false,
		},
		{
			name: "regex match",
			mapping: &Mapping{
				Match: "f.*",
			},
			path:     ".config/bar/foo",
			expected: true,
		},
		{
			name: "regex mismatch",
			mapping: &Mapping{
				Match: "f.*",
			},
			path:     "bar",
			expected: false,
		},
		{
			name: "target os match",
			mapping: &Mapping{
				Match:    "foo",
				TargetOS: runtime.GOOS,
			},
			path:     "foo",
			expected: true,
		},
		{
			name: "target os mismatch",
			mapping: &Mapping{
				Match:    "foo",
				TargetOS: notCurrentOS(),
			},
			path:     "foo",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mapping.IsMatch(tc.path) != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, !tc.expected)
			}
		})
	}
}

func TestAction_String(t *testing.T) {
	cases := []struct {
		name     string
		action   Action
		expected string
	}{
		{
			name:     "none",
			action:   ActionNone,
			expected: "NONE",
		},
		{
			name:     "skip",
			action:   ActionSkip,
			expected: "SKIP",
		},
		{
			name:     "translate",
			action:   ActionTranslate,
			expected: "TRANSLATE",
		},
		{
			name:     "link as dir",
			action:   ActionLinkAsDir,
			expected: "LINK_AS_DIR",
		},
		{
			name:     "unknown",
			action:   Action(99),
			expected: "UNKNOWN",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.action.String() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.action.String())
			}
		})
	}
}

func TestAction(t *testing.T) {
	cases := []struct {
		name     string
		mapping  *Mapping
		expected Action
	}{
		{
			name:     "none",
			mapping:  &Mapping{},
			expected: ActionNone,
		},
		{
			name: "skip",
			mapping: &Mapping{
				Skip: true,
			},
			expected: ActionSkip,
		},
		{
			name: "link as dir",
			mapping: &Mapping{
				LinkAsDir: true,
			},
			expected: ActionLinkAsDir,
		},
		{
			name: "translate",
			mapping: &Mapping{
				Dest: "bar",
			},
			expected: ActionTranslate,
		},
		{
			name: "skip overrules link as dir",
			mapping: &Mapping{
				Skip:      true,
				LinkAsDir: true,
			},
			expected: ActionSkip,
		},
		{
			name: "skip overrules translate",
			mapping: &Mapping{
				Skip: true,
				Dest: "bar",
			},
			expected: ActionSkip,
		},
		{
			name: "link as dir overrules translate",
			mapping: &Mapping{
				LinkAsDir: true,
				Dest:      "bar",
			},
			expected: ActionLinkAsDir,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mapping.Action() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.mapping.Action())
			}
		})
	}
}
