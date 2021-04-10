"""
Usage: dfm sync [options]

Options:
    -m, --message <msg>   Use the given <msg> as the commit message.
    -n, --name <name>     Sync only the given profile or module by name.
    -d, --dry-run         Print git commands instead of executing them.

Sync the current profile and modules.
"""

from os.path import exists
from sys import exit

from dfm.cli.utils import load_profile, profile_dir
from dfm.profile import Profile


def find_module(name, profile):
    """Recursively search profile and it's modules for a module with name."""
    for mod in profile.modules:
        if mod.name == name:
            return mod

        if mod.modules:
            found = find_module(name, mod)
            if found is not None:
                return found

    return None


def run(args):
    """Run profile.sync for the current profile."""
    if args["--name"]:
        possible_dir = profile_dir(args["--name"])
        if not exists(possible_dir):
            curprofile = load_profile()
            profile = find_module(args["--name"], curprofile)
            if profile is None:
                print("no module or profile matched name: {}".format(args["--name"]))
                exit(1)
        else:
            profile = Profile.load(possible_dir)
    else:
        profile = load_profile()

    profile.sync(
        dry_run=args["--dry-run"],
        commit_msg=args["--message"],
    )
