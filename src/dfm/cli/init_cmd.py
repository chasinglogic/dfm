"""
Usage: dfm init [options] <name>

Create a profile with the given name.
"""

import os
import sys

from dfm.config import dfm_dir
from dfm.profile import Profile


def run(args):
    """Create a new Profile with the given name."""
    new_profile_dir = os.path.join(dfm_dir(), "profiles", args["<name>"])

    try:
        os.makedirs(new_profile_dir)
    except OSError as exc:
        print("Failed to create profile directory!")
        print("Error:", exc)
        sys.exit(1)

    Profile.new(new_profile_dir)
