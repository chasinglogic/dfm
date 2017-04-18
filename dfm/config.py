"""
Contains the config functions and global for dfm

dfm, a dotfile manager for lazy people and pair programmers

Copyright (C) 2016 Mathew Robinson <mathew.robinson3114@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
"""

import os
import json


def get_default_config_dir():
    """Get the default config dir for dfm."""
    xdg = os.environ.get("XDG_CONFIG_HOME", "")
    if xdg == "":
        home = os.environ.get("HOME", "")
        return os.path.join(home, ".config", "dfm")
    return os.path.join(xdg, "dfm")


def load_config(config_dir):
    """Load config from the given config_dir."""
    global CONFIG_FILE
    CONFIG_FILE = os.path.join(config_dir, "config.json")
    if os.path.isfile(CONFIG_FILE):
        with open(CONFIG_FILE, "r") as cf:
            return json.load(cf)
    return {}


def save_config():
    """Save CONFIG to the appropriate location."""
    with open(CONFIG_FILE, "w") as cf:
        json.dump(CONFIG, cf)


def upgrade_config(old_conf):
    """Upgrade the old style config to the new style."""
    with open(old_conf) as f:
        jsn = json.load(f)
    if jsn is None or jsn == {}:
        print("Old config was empty nothing to do.")
        return
    CONFIG['verbose'] = jsn['Verbose']
    CONFIG['profile'] = os.path.join(CONFIG_DIR, jsn['CurrentProfile'])
    save_config()


CONFIG_DIR = get_default_config_dir()
CONFIG = load_config(CONFIG_DIR)
