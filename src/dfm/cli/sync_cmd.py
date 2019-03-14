"""
Usage: dfm sync [options]

Options:
    -m, --message <msg>  Use the given <msg> as the commit message.

Sync the current profile and modules.
"""

from dfm.cli.utils import inject_profile


@inject_profile
def run(args, profile=None):
    """Run profile.sync for the current profile."""
    if args['--message']:
        profile.commit_msg = args['--message']
    profile.sync()
