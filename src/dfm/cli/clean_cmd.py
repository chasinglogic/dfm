"""
Usage: dfm clean

Removes broken symlinks. This can clean up a cluttered $HOME directory after
you've removed dotfiles from your profile.
"""

import logging
import os
import re
import shutil
import textwrap

from yaspin import yaspin

from dfm.cli.utils import inject_profile
from dfm.config import xdg_dir

logger = logging.getLogger(__name__)

ANSI_ESCAPE = re.compile(r"(\x9B|\x1B\[)[0-?]*[ -\/]*[@-~]")


def clean_links(directory, profile_dir):
    """Remove all broken symlinks in directory."""
    max_width, _ = shutil.get_terminal_size()
    max_width -= 4
    with yaspin() as spinner:
        for dirpath, _, files in os.walk(directory):
            msg = "Scanning for dead links in {}".format(
                ANSI_ESCAPE.sub("", dirpath),
            )
            spinner.text = textwrap.shorten(msg, max_width)
            for file in files:
                ab_path = os.path.join(dirpath, file)
                if not os.path.islink(ab_path):
                    continue

                path = os.readlink(ab_path)
                if profile_dir not in path:
                    logger.debug("Skipping non-profile dead link: %s", ab_path)
                    continue

                if not os.path.exists(path):
                    logger.info("Removing dead link: %s", ab_path)
                    os.unlink(ab_path)


@inject_profile
def run(_args, profile):
    """Run the clean subcommand."""
    home = os.getenv("HOME")
    xdg = xdg_dir()
    if home:
        clean_links(home, profile.link_manager.where)
    clean_links(xdg, profile.link_manager.where)
