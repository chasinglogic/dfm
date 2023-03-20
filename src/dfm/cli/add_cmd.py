"""
Usage: dfm add [options] <file>...

Add files to the current dotfile profile doing "reverse dotfile-ization" on them
and linking back correctly.
"""

import os
import shutil
import sys

from dfm.cli.utils import inject_profile


@inject_profile
def run(args, profile):
    """Run the add command with args."""
    for filename in args["<file>"]:
        oldfile = os.path.abspath(filename)
        if not os.path.exists(oldfile):
            print("error: {}: file does not exist".format(oldfile), file=sys.stderr)
            sys.exit(1)

        newfile = os.path.relpath(oldfile, profile.link_manager.target_dir)
        newfile = os.path.join(profile.link_manager.where, newfile)
        parent = os.path.dirname(newfile)
        if not os.path.isdir(parent):
            os.makedirs(parent)

        if os.path.isfile(oldfile):
            shutil.copy2(oldfile, newfile)
            os.remove(oldfile)
        else:
            shutil.copytree(oldfile, newfile)
            shutil.rmtree(oldfile)

    print("I should not be syncing.")
    # profile.sync(skip_modules=True)
    profile.link()
