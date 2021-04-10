"""
Usage: dfm remove <profile>

Lists currently downloaded and available profiles.
"""

import os
import shutil
import sys

from dfm.config import dfm_dir


def run(args):
    """Run the remove subcommand."""
    profile_p = os.path.join(dfm_dir(), "profiles", args["<profile>"])
    if not os.path.isdir(profile_p):
        print("no profile with that name exists")
        sys.exit(1)

    ans = input("Remove {}? [Y/n]:".format(profile_p))
    if ans.lower().startswith("n"):
        return

    shutil.rmtree(profile_p)
