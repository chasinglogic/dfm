"""
dfm, a dotfile manager for lazy people and pair programmers.

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

import click
import os
import sys
from shutil import rmtree, which

from dfm.lib import pull_profile
from dfm.lib import push_profile
from dfm.lib import clone_profile
from dfm.lib import link_profile
from dfm.lib import checkout_profile
from dfm.lib import create_and_init_profile
from dfm.lib import commit_profile
from dfm.lib import set_remote_profile

from dfm.lib import add_file
from dfm.lib import get_repo_url
from dfm.lib import get_profile_path
from dfm.lib import git_pass_through

from dfm.config import CONFIG
from dfm.config import CONFIG_DIR
from dfm.config import get_default_config_dir
from dfm.config import load_config
from dfm.config import save_config
from dfm.config import upgrade_config


LICENSE = """
dfm, a dotfile manager for lazy people and pair programmers

Copyright (C) 2016 Mathew Robinson <chasinglogic@gmail.com>

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

__version__ = '2.1'


@click.group()
@click.option('--verbose', '-vv', is_flag=True)
@click.option('--config', '-c',
              help='The path where dfm stores it\'s config and profiles.',
              default=get_default_config_dir(),
              type=click.Path(resolve_path=True))
def dfm(verbose, config):
    """A dotfile manager for lazy people and pair programmers."""
    if config != get_default_config_dir():
        CONFIG = load_config(config)

    # If CONFIG_DIR does not exist, make it so #1
    if not os.path.isdir(CONFIG_DIR):
        os.mkdir(CONFIG_DIR)

    if which('git') is None:
        print('Git is not in the $PATH.')
        print('Git is required for dfm please install then try again.')
        sys.exit(1)

    if verbose:
        CONFIG['verbose'] = True

    pass


@dfm.command()
@click.option('--branch', '-b',
              help='Branch you would like to pull from.',
              default='master')
def pull(branch):
    """Pull changes from the remote."""
    profile = CONFIG.get('profile', None)
    click.echo('Updating profile %s' % profile)
    if profile:
        pull_profile(profile, branch)


@dfm.command()
@click.option('--branch', '-b',
              help='Branch you would like to pull from.',
              default='master')
def push(branch):
    """Push local changes to the remote."""
    profile = CONFIG.get('profile', None)
    click.echo('Pushing profile %s' % profile)
    if profile:
        push_profile(profile, branch)


@dfm.command(context_settings={'ignore_unknown_options': True})
@click.argument('args', nargs=-1)
def commit(args):
    """Run a git commit for the current profile."""
    profile = CONFIG.get('profile')
    commit_profile(profile, list(args))


@dfm.command()
@click.option('--link', '-l', is_flag=True,
              help='Link the profile after downloading it.')
@click.option('--force', '-f', is_flag=True,
              help='Force removal of non-symlink type files')
@click.argument('repo')
def clone(link, force, repo):
    """Clone a profile from a git repo."""
    repo_url = get_repo_url(repo)
    profile_path = get_profile_path(CONFIG_DIR, repo)
    click.echo('Creating profile %s from %s' % (profile_path, repo_url))
    clone_profile(repo_url, profile_path)
    if link:
        link_profile(profile_path, force)
        CONFIG['profile'] = profile_path
        save_config()


@dfm.command()
@click.option('--force', '-f', is_flag=True,
              help='Force removal of non-symlink type files')
@click.argument('profile')
def link(force, profile):
    """Link the profile with the given name."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    link_profile(profile_path, force)
    CONFIG['profile'] = profile_path
    save_config()


@dfm.command()
@click.argument('profile')
def init(profile):
    """Create an empty profile with the given name."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    create_and_init_profile(profile_path)


@dfm.command()
@click.argument('profile')
def rm(profile):
    """Remove the profile with the given name."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    click.echo('Removing profile %s' % profile_path)
    rmtree(profile_path)


@dfm.command()
@click.argument('path', nargs=-1,
                type=click.Path(resolve_path=True, exists=True))
def add(path):
    """Add a file or directory to the current profile."""
    profile = CONFIG.get('profile')
    for f in path:
        add_file(f, profile)


@dfm.command()
@click.argument('branch')
def checkout(branch):
    """Switch to a different branch for the active profile."""
    profile = CONFIG.get('profile', None)
    if profile:
        checkout_profile(profile, branch)
    else:
        click.echo('No profile currently active.')


@dfm.command()
def version():
    """Show the current dfm version."""
    print('You are running dfm version %s' % __version__)


@dfm.command()
def license():
    """Show dfm licensing info."""
    print('\nYou are running dfm version %s' % __version__)
    print(LICENSE)


@dfm.command()
@click.argument('remote')
def remote(remote):
    """Set the git remote for the current profile."""
    profile = CONFIG.get('profile', None)
    if profile:
        set_remote_profile(profile, remote)


@dfm.command(context_settings={'ignore_unknown_options': True})
@click.argument('args', nargs=-1)
def git(args):
    """Run the given git command in the current profile."""
    profile = CONFIG.get('profile')
    if profile:
        git_pass_through(profile, ['git'] + list(args))
    else:
        print('No git command specified.')


@dfm.command()
def where():
    """Return the path to the current profile. Useful for piping."""
    print(CONFIG.get('profile'))


@dfm.command()
@click.option('--config', '-c',
              help='Path to the old config',
              default=os.path.join(os.getenv('HOME'), '.dfm'))
def upgrade(config):
    """
    Upgrade from the old style config to the new style.

    If you're switching from the go version to the python version you
    should run this command. Should only be run once.
    """
    if os.path.exists(config):
        upgrade_config(config)
        print('Config upgraded, would you like to remove the old config?')
        ans = input('y/N: ')
        if 'y' in ans:
            print('Removing:', config)
            os.remove(config)
        return
    print('No old config found you\'re good to go!')
