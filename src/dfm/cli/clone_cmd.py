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

import os
import subprocess
import sys

from dfm.cli.link_cmd import run as run_link
from dfm.profile import dfm_dir, get_name


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
