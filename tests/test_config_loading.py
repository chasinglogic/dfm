"""Test dotfile config loading"""

import os
from tempfile import TemporaryDirectory

from dfm.profile import Profile


def setup_module():
    """Set up the DFM_CONFIG_DIR for this module run."""
    os.environ["DFM_CONFIG_DIR"] = TemporaryDirectory().name
    os.environ["DFM_DISABLE_MODULES"] = "1"


def test_list_files(dotfile_dir):
    """Test that a profile properly lists it's directory."""
    dotfiles, directory = dotfile_dir()
    profile = Profile(directory)
    profile._find_files()
    expected_files = sorted(dotfiles)
    assert sorted(profile.files) == expected_files
    assert not profile.always_sync_modules


def test_config_loading(dotdfm):
    """Test that a profile properly loads the config file."""
    _, directory = dotdfm(
        """
always_sync_modules: true
"""
    )
    profile = Profile(directory)
    assert profile.always_sync_modules


def test_mapping_loading(dotdfm):
    """Test that a profile properly loads the config file mappings."""
    _, directory = dotdfm(
        """
mappings:
    - match: emacs
      skip: true
    - match: vimrc
      dest: .config/nvim/init.vim
"""
    )
    profile = Profile(directory)
    assert len(profile.mappings) == 8


def test_module_loading(dotdfm):
    """Test that a profile properly loads the config file modules."""
    _, directory = dotdfm(
        """
modules:
  - repo: https://github.com/robbyrussell/oh-my-zsh
    link: none
    location: ~/.oh-my-zsh
"""
    )
    profile = Profile(directory)
    assert profile.modules
