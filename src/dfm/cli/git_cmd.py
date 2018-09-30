"""
Usage: dfm git [<args>....]

Run the git command in the current dotfile profile.
"""

import subprocess
import sys

from dfm.cli.utils import current_profile


def run(args):
    """Run the git subcommand with args."""
    proc = subprocess.Popen(
        ['git'] + args['<args>'],
        cwd=current_profile(),
        stdin=sys.stdin,
        stdout=sys.stdout,
        stderr=sys.stderr)
    proc.wait()
