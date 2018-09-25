"""Usage: dfm git [<args>....]

Run the git command in the current dotfile profile.
"""

import sys
import subprocess

from dfm.cli.utils import current_profile


def run(args):
    proc = subprocess.Popen(
        ['git'] + args['<args>'],
        cwd=current_profile(),
        stdin=sys.stdin,
        stdout=sys.stdout,
        stderr=sys.stderr)
    proc.wait()
