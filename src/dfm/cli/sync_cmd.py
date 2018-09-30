"""Usage:
  dfm sync [options]

Sync the current profile and modules.
"""

from dfm.cli.utils import inject_profile


@inject_profile
def run(_args, profile=None):
    profile.sync()
