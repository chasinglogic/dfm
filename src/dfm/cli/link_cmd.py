"""
Usage:
  dfm link [options] [<profile>]

Link / activate dfm dotfile profiles. If no profile is provided relinks the
current profile.

Options:
    -d, --dry-run    If provided simply print what links would be generated
    -o, --overwrite  If provided dfm will delete files and directories which
                     exist at the target link locations. DO NOT USE THIS IF YOU
                     ARE UNSURE AS IT WILL RESULT IN DATA LOSS.
"""

import logging

from dfm.cli.utils import current_profile, load_profile, switch_profile


def run(args):
    """Run the link subcommand, setting the current_profile."""
    dry_run = args['--dry-run']
    if args['<profile>'] and not dry_run:
        profile = switch_profile(args['<profile>'])
    elif args['<profile>']:
        profile = load_profile(args['<profile>'])
    else:
        profile = load_profile(current_profile())

    links = profile.link(overwrite=args['--overwrite'], dry_run=dry_run)
