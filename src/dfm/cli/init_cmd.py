"""
Usage: dfm init [options] <name>

Create a profile with the given name.
"""

import os
from dfm.dotfile import DotfileRepo, dfm_dir


def run(args):
    new_profile_dir = os.path.join(dfm_dir(), "profiles", args["<name>"])

    try:
        os.makedirs(new_profile_dir)
    except Exception as e:
        print("Failed to create profile directory!")
        print("Error:", e)
        os.exit(1)

    profile = DotfileRepo(new_profile_dir)
    profile.init()
