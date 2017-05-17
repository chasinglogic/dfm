"""Contains all of the business logic for dfm."""

import os
import shutil

from pathlib import Path
from subprocess import run

from dfm.config import CONFIG


def get_repo_url(repo_name):
    """
    Determine if we need to add github.com or not.

    returns appropriate url.
    """
    if len(repo_name.split('/')) == 2:
        return 'https://github.com/' + repo_name
    return repo_name


def get_profile_path(config_dir, profile_name):
    """Get the appropriate dir path for profile_name."""
    spl = profile_name.split('/')
    if len(spl) > 1:
        # if we got passed a url of some sort return
        # the second to last element of the split.
        return os.path.join(config_dir, 'profiles', spl[len(spl) - 2])
    return os.path.join(config_dir, 'profiles', profile_name)


def pull_profile(path, branch):
    """Run git pull at path for branch."""
    run(['git', 'pull', 'origin', branch], cwd=path)


def clone_profile(repo, profile_path):
    """Run git clone for repo, cloning into profile_path."""
    run(['git', 'clone', repo, profile_path])


def push_profile(path, branch):
    """Run git push origin for branch at path."""
    run(['git', 'push', 'origin', branch], cwd=path)


def create_and_init_profile(profile_path):
    """Create a profile at profile_path and run git init."""
    print('Creating profile %s' % profile_path)
    os.mkdir(profile_path)
    run(['git', 'init'], cwd=profile_path)


def gen_dot_file(filename):
    """Take filename and generate the appropriate new path for it."""
    df = filename if filename.startswith('.') else '.' + filename
    return os.path.abspath(os.path.join(os.environ.get('HOME', ''), df))


def link_file(fle, dotfile, force=False):
    """Create a symlink from fle to dotfile if appropriate."""
    # Check if a non sym linked version exists.
    if ((os.path.exists(dotfile) and not os.path.islink(dotfile))
       and not force):
        print('Error linking: %s' % dotfile)
        print('Dotfile exists and isn\'t symlink.')
        print('Refusing to overwrite. Use --force to overwrite')
        return

    if os.path.isfile(dotfile) or os.path.islink(dotfile):
        os.remove(dotfile)
    if os.path.isdir(dotfile):
        shutil.rmtree(dotfile)

    os.symlink(fle, dotfile)
    if CONFIG.get('verbose', False):
        print('Linked file %s -> %s' % (fle, dotfile))


def xdg_default():
    """Return the default $XDG_CONFIG_HOME if not set."""
    return os.path.join(os.environ.get('HOME', '/'), '.config')


def link_profile(profile_path, force=False):
    """Link all dotfiles in profile_path."""
    abpath = os.path.abspath(profile_path)
    dot_files = os.scandir(abpath)
    print('Linking profile %s' % profile_path)
    for d in dot_files:
        # Skip the git directory and dfm config file
        if (
                d.name == '.git'
                or d.name == '.dfm.yml'
                or d.name == '.dfm.yaml'
        ):
            continue

        # If files are in the config dir, that means they need to
        # follow $XDG_CONFIG_HOME
        if ((d.name == 'config' or d.name == '.config')
           and os.path.isdir(os.path.join(abpath, d.name))):
            # Get the .config directory files
            dfgf = os.scandir(d.path)
            for f in dfgf:
                # .config files have a different path
                xdg = os.environ.get('XDG_CONFIG_HOME', xdg_default())
                dfp = os.path.join(xdg, f.name)
                link_file(f.path, dfp, force=force)
            continue

        # otherwise just link it
        link_file(d.path, gen_dot_file(d.name), force=force)


def add_file(path, profile):
    """Add file to the current profile."""
    xdg = os.environ.get('XDG_CONFIG_HOME', '')
    if os.path.islink(path):
        print('You can\'t add a symlink using dfm.')
        return

    if not os.path.exists(path):
        print('No such file with name:', path)
        return

    old_file = Path(path)
    fn = old_file.name
    # If starts with a dot remove it.
    if fn.startswith('.'):
        fn = fn.replace('.', '', 1)

    new_path = os.path.join(profile, fn)
    if Path(xdg) in {old_file}:
        new_path = os.path.join(profile, 'config', fn)

    os.rename(bytes(old_file), new_path)
    link_file(new_path, bytes(old_file), force=True)
    run(['git', 'add', fn], cwd=profile)


def checkout_profile(profile, branch):
    """Checkout the git branch on profile."""
    run(['git', 'checkout', branch], cwd=profile)


def commit_profile(profile, args):
    """Run git commit with the given message."""
    run(['git', 'commit'] + args, cwd=profile)


def set_remote_profile(profile, remote):
    """Set origin for profile to remote."""
    run(['git', 'remote', 'set-url', 'origin', remote], cwd=profile)


def git_pass_through(profile, argv):
    """Run the specified git command with cwd=profile."""
    run(argv, cwd=profile)
