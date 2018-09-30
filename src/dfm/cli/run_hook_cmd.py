"""
Usage: dfm run_hook <hook>

Runs <hook> without the need to invoke the side effects of the given action.
"""

from dfm.cli.utils import inject_profile


@inject_profile
def run(args, profile=None):
    """Run hook with the given name in the .dfm.yml."""
    profile.run_hook(args['run_hook'])
