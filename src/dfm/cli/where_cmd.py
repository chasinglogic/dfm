"""
Usage: dfm where

Prints the location of the current profile
"""

from dfm.cli.utils import current_profile


def run(_args):
    """Print current_profile path."""
    print(current_profile())
