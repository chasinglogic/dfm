"""
Usage: dfm git [<args>....]

Run the git command in the current dotfile profile.
"""

from dfm.cli.utils import current_profile


def run(args):
    """Run the git subcommand with args."""
    profile = current_profile()
    profile.df_repo.git(args["<args>"])
