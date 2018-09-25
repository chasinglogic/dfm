"""Usage: dfm clone [options] <url>

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

import subprocess
import sys
import os

from dfm.cli.link_cmd import run as run_link
from dfm.dotfile import dfm_dir


def get_name(url):
    return url.split('/')[-1]


def run(args):
    name = args.get('--name', get_name(args['<url>']))
    path = os.path.join(dfm_dir(), 'profiles', name)
    subprocess.call(['git', 'clone', args['<url>'], path])
    if args['--link']:
        args = {
            '<profile>': name,
            '--overwrite': args['--overwrite'],
            '--dry-run': False,
        }

        run_link(args)
