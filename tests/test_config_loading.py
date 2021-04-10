"""Test dotfile config loading"""

import os
from tempfile import TemporaryDirectory

from dfm.profile import Profile


def setup_module():
    """Set up the DFM_CONFIG_DIR for this module run."""
    os.environ["DFM_CONFIG_DIR"] = TemporaryDirectory().name
    os.environ["DFM_DISABLE_MODULES"] = "1"


def test_config_loading(dotdfm):
    """Test that a profile properly loads the config file."""
    _, directory = dotdfm(
        """
pull_only: true
""",
    )
    profile = Profile.load(directory)
    assert profile.pull_only


def test_mapping_loading(dotdfm):
    """Test that a profile properly loads the config file mappings."""
    _, directory = dotdfm(
        """
mappings:
    - match: emacs
      skip: true
    - match: vimrc
      dest: .config/nvim/init.vim
""",
    )
    profile = Profile.load(directory)
    assert len(profile.link_manager.mappings) == 8


def test_module_loading(dotdfm):
    """Test that a profile properly loads the config file modules."""
    _, directory = dotdfm(
        """
modules:
  - repo: https://github.com/robbyrussell/oh-my-zsh
    link: none
    location: ~/.oh-my-zsh
""",
    )
    profile = Profile.load(directory)
    assert profile.modules
