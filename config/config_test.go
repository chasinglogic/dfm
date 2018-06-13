package config

import "testing"

func TestConfigMethods(t *testing.T) {
	cfg := Config{
		Dir:                "/home/chasinglogic/.config/dfm",
		CurrentProfileName: "chasinglgoic",
	}

	expected := "/home/chasinglogic/.config/dfm/profiles"
	if cfg.ProfileDir() != expected {
		t.Errorf("Expected %s Got %s", expected, cfg.ProfileDir())
	}

	expected = "/home/chasinglogic/.config/dfm/profiles/chasinglogic"
	if cfg.CurrentProfile() != expected {
		t.Errorf("Expected %s Got %s", expected, cfg.CurrentProfile())
	}

	expected = "/home/chasinglogic/.config/dfm/profiles/lionize"
	if cfg.GetProfileByName("lionize") != expected {
		t.Errorf("Expected %s Got %s", expected, cfg.GetProfileByName("lionize"))
	}
}
