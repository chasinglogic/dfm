// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.


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
