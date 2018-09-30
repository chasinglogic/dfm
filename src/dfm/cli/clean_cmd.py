"""
Usage: dfm clean

Removes broken symlinks. This can clean up a cluttered $HOME directory after
you've removed dotfiles from your profile.
"""

import os

from dfm.dotfile import xdg_dir


def clean_links(directory):
    """Remove all broken symlinks in directory."""
    for filename in os.listdir(directory):
        ab_path = os.path.join(directory, filename)
        if os.path.islink(ab_path) and not os.path.exists(ab_path):
            print('Removing dead link:', ab_path)
            os.unlink(ab_path)


def run(_args):
    """Run the clean subcommand."""
    home = os.getenv('HOME')
    xdg = xdg_dir()
    if home:
        clean_links(home)
    clean_links(xdg)
