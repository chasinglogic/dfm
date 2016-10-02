#
#     dfm, a dotfile manager for lazy people and pair programmers
#
#     Copyright (C) 2016 Mathew Robinson <mathew.robinson3114@gmail.com>
#
#     This program is free software: you can redistribute it and/or modify
#     it under the terms of the GNU General Public License as published by
#     the Free Software Foundation, either version 3 of the License, or
#     (at your option) any later version.
#
#     This program is distributed in the hope that it will be useful,
#     but WITHOUT ANY WARRANTY; without even the implied warranty of
#     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#     GNU General Public License for more details.
#
#     You should have received a copy of the GNU General Public License
#     along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
import click
import os
import sys
from shutil import rmtree, which
from dfm.lib import *
from dfm.config import *


LICENSE = """
dfm, a dotfile manager for lazy people and pair programmers

Copyright (C) 2016 Mathew Robinson <mathew.robinson3114@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

"""
VERSION_NUMBER = 0.3

@click.group()
@click.option("--verbose", "-vv",
              default=False,
              is_flag=True)
@click.option("--config", "-c", 
              help="The path where dfm stores it's config and profiles.",
              default=get_default_config_dir(),
              type=click.Path(resolve_path=True))
def dfm(verbose, config):
    """A dotfile manager for lazy people and pair programmers."""
    if config != get_default_config_dir():
        set_config(config)

    # If CONFIG_DIR does not exist, make it so #1
    if not os.path.isdir(CONFIG_DIR):
        os.mkdir(CONFIG_DIR)

    if which("git") == None:
        click.echo("Git is not in the $PATH. Git is required for dfm please install then try again.")
        sys.exit(1)

    if verbose:
        CONFIG["verbose"] = True

    pass

@dfm.command()
@click.option("--branch", "-b",
              help="Branch you would like to pull from.",
              default="master")
def pull(branch):
    """Pull changes from the remote."""
    profile = CONFIG.get("profile", None)
    click.echo("Updating profile %s" % profile)
    if profile:
        pull_profile(profile, branch)

@dfm.command()
@click.option("--branch", "-b",
              help="Branch you would like to pull from.",
              default="master")
def push(branch):
    """Push local changes to the remote."""
    profile = CONFIG.get("profile", None)
    click.echo("Pushing profile %s" % profile)
    if profile:
        push_profile(profile, branch)

@dfm.command()
@click.argument("message")
def commit(message):
    """Run a git commit for the current profile."""
    profile = CONFIG.get("profile")
    commit_profile(profile, message)

@dfm.command()
@click.option("--link", "-l",
              default=False,
              is_flag=True,
              help="Link the profile after downloading it.")
@click.option("--force", "-f",
              default=False,
              is_flag=True,
              help="Force removal of non-symlink type files")
@click.argument("repo")
def clone(link, force, repo):
    """Clone a profile from a git repo."""
    repo_url = get_repo_url(repo)
    profile_path = get_profile_path(CONFIG_DIR, repo)
    click.echo("Creating profile %s from %s" % (profile_path, repo_url))
    clone_profile(repo_url, profile_path)
    if link:
        link_profile(profile_path, force)
        CONFIG["profile"] = profile_path
        save_config()

@dfm.command()
@click.option("--force", "-f",
              help="Force removal of non-symlink type files",
              is_flag=True,
              default=False)
@click.argument("profile")
def link(force, profile):
    """Link the profile with the given name."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    link_profile(profile_path, force)
    CONFIG["profile"] = profile_path
    save_config()

@dfm.command()
@click.argument("profile")
def init(profile):
    """Create an empty profile with the given name."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    create_and_init_profile(profile_path)

@dfm.command()
@click.argument("profile")
def rm(profile):
    """Remove the profile with the given name."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    click.echo("Removing profile %s" % profile_path)
    rmtree(profile_path)

@dfm.command()
@click.argument("path",
                type=click.Path(resolve_path=True, exists=True),
                nargs=-1)
def add(path):
    """Add a file or directory to the current profile."""
    profile = CONFIG.get("profile")
    for f in path:
        add_file(f, profile)

@dfm.command()
@click.argument("branch")
def chk(branch):
    """Switch to a different branch for the active profile."""
    profile = CONFIG.get("profile", None)
    if profile:
        checkout_profile(profile, branch)
    else:
        click.echo("No profile currently active.")

@dfm.command()
def version():
    """Show the current dfm version."""
    print("You are running dfm version %.1f" % VERSION_NUMBER)

@dfm.command()
def license():
    """Show dfm licensing info."""
    print("\nYou are running dfm version %.1f" % VERSION_NUMBER)
    print(LICENSE)

@dfm.command()
@click.argument("remote")
def remote(remote):
    """Set the git remote for the current profile."""
    profile = CONFIG.get("profile", None)
    if profile:
        set_remote_profile(profile, remote)
