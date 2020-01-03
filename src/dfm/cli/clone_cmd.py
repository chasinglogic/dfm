"""
Usage: dfm clone [options] <url>

Options:
    -n <name>, --name <name>  Defaults to the 'basename' of the url. For
                              example: https://github.com/chasinglogic/dotfiles
                              would be 'dotfiles'
    -l, --link                If provided the profile will be immediately linked
    -o, --overwrite           If provided links will overwrite files and
                              directories that exist at target locations. DO NOT
                              USE THIS IF YOU ARE UNSURE SINCE IT WILL RESULT IN
                              DATA LOSS.

Clones the repository at <url> to a new profile with <name>
"""

import sys
import os
import subprocess

from dfm.cli.link_cmd import run as run_link
from dfm.dotfile import dfm_dir


def get_name(url):
    """
    Generate a profile name based on the git url.

    This directly corresponds to the second to last element in the URL.  For
    example: https://github.com/chasinglogic/dotfiles would be 'chasinglogic'

    In the case of an ssh url or other url will correspond to the
    first path element. For example:
    git@github.com:chasinglogic/dotfiles would be 'chasinglogic'
    """
    try:
        if url.find("@") > 1:
            return url.split(":")[-1].split("/")[0]
        else:
            return url.split("/")[-2]
    # Any kind of exception i.e. IndexError or failure to split we
    # just return nothing.
    except Exception:
        return ""


def run(args):
    """Run the clone command."""
    name = args["--name"]
    if name is None:
        name = get_name(args["<url>"])

    if not name:
        print(
            "--name flag not provided and could not be determined from",
            args["<url>"],
            "please provide a name for this profile via the --name flag",
        )
        sys.exit(1)

    path = os.path.join(dfm_dir(), "profiles", name)
    subprocess.call(["git", "clone", args["<url>"], path])
    if args["--link"]:
        args = {
            "<profile>": name,
            "--overwrite": args["--overwrite"],
            "--dry-run": False,
        }

        run_link(args)
