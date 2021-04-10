"""
Usage: dfm list

Lists currently downloaded and available profiles.
"""

import os

from dfm.config import dfm_dir


def run(_args):
    """List all available profiles on this system."""
    profiles_dir = os.path.join(dfm_dir(), "profiles")
    if not os.path.isdir(profiles_dir):
        print("There are no profiles on this system yet. Create one with `dfm init`!")
        print(
            "For more information see the dfm documentation: "
            "https://github.com/chasinglogic/dfm",
        )
        return

    for profile in os.listdir(profiles_dir):
        if not profile.startswith("."):
            print(profile)
