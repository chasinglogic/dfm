"""Test dotfile config loading"""

import os

from tempfile import TemporaryDirectory
from dfm.dotfile import Profile


def setup_module():
    """ setup any state specific to the execution of the given module."""
    os.environ['DFM_CONFIG_DIR'] = TemporaryDirectory().name


def test_list_files(dotfile_dir):
    dotfiles, directory = dotfile_dir
    profile = Profile(directory)
    assert sorted(profile.files) == sorted(dotfiles)
    assert not profile.always_sync_modules


def test_config_loading(dotdfm):
    _, directory = dotdfm("""
always_sync_modules: true
""")
    profile = Profile(directory)
    assert profile.always_sync_modules


def test_mapping_loading(dotdfm):
    _, directory = dotdfm("""
mappings:
    - match: emacs
      skip: true
    - match: vimrc
      location: ~/.config/nvim/init.vim
""")
    profile = Profile(directory)
    assert len(profile.mappings) == 10


def test_module_loading(dotdfm):
    _, directory = dotdfm("""
modules:
  - repo: https://github.com/robbyrussell/oh-my-zsh
    link: none
    location: ~/.oh-my-zsh
  - repo: keybase://private/chasinglogic/secrets
  - repo: keybase://private/chasinglogic/personal-infrastructure
    link: none
  - repo: keybase://private/chasinglogic/Notes
    location: ~/Notes
  - repo: https://github.com/tmux-plugins/tpm
    link: none
    pull_only: true
    location: ~/.tmux/plugins/tpm
""")
    profile = Profile(directory)
    assert profile.modules
