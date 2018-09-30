"""
Usage: dfm add [options] <file>...

Add files to the current dotfile profile doing "reverse dotfile-ization" on them
and linking back correctly.

Options:
    -k, --keep-dot  If provided the file will be copied exactly as is. Use this
                    if you keep your dotfiles repo files as actual dotfiles. For
                    convenience you can set the environment variable
                    $DFM_KEEP_DOT and omit this flag.
"""

import os
import shutil
import sys

from dfm.cli.utils import inject_profile


@inject_profile
def run(args, profile=None):
    """Run the add command with args."""
    for filename in args['<file>']:
        oldfile = os.path.abspath(filename)
        if not os.path.exists(oldfile):
            print('{}: file does not exist'.format(oldfile))
            sys.exit(1)

        if args['--keep-dot'] or not oldfile.startswith('.'):
            newfile = os.path.basename(oldfile)
        else:
            # Stip off the leading dot
            newfile = os.path.basename(oldfile)[1:]

        newfile = os.path.join(profile.where, newfile)
        if os.path.isfile(oldfile):
            shutil.copy2(oldfile, newfile)
            os.remove(oldfile)
        else:
            shutil.copytree(oldfile, newfile)
            shutil.rmtree(oldfile)

    profile.sync(skip_modules=True)
    profile.link()
