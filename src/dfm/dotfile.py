"""Translate file names to the appropriate targets."""

import logging
import os
import re
import shlex
import subprocess
import sys
from shutil import rmtree

import yaml

logger = logging.getLogger(__name__)


def xdg_dir():
    """Return the XDG_CONFIG_HOME or default."""
    if os.getenv('XDG_CONFIG_HOME'):
        return os.getenv('XDG_CONFIG_HOME')
    return os.path.join(os.getenv('HOME'), '.config')


def dfm_dir():
    """Return the dfm configuration / state directory."""
    if os.getenv('DFM_CONFIG_DIR'):
        return os.getenv('DFM_CONFIG_DIR')
    return os.path.join(xdg_dir(), 'dfm')


class Mapping:
    """
    Maps a filename to a new destination.

    Allows for dotfiles to be skipped, or redirected to a target directory other than 'HOME'
    """

    def __init__(self, match, target_dir='', skip=False):
        self.match = match
        self.target_dir = target_dir.replace('~', os.getenv('HOME'))
        self.skip = skip
        self.rgx = re.compile(match)

    @classmethod
    def from_dict(cls, config):
        """Return a Mapping from the config dictionary"""
        return cls(**config)

    def matches(self, path):
        """Determine if this mapping matches path."""
        return self.rgx.search(path)


DEFAULT_MAPPINGS = [
    Mapping(
        r'\/\.git\/',
        skip=True,
    ),
    Mapping(
        r'\/.gitignore$',
        skip=True,
    ),
    Mapping(
        r'\/LICENSE(\.md)?$',
        skip=True,
    ),
    Mapping(
        r'\/\.dfm\.yml$',
        skip=True,
    ),
    Mapping(
        r'\/README(\.md)?$',
        skip=True,
    ),
]


def unable_to_remove(filename, overwrite=False):
    """Remove the file if necessary. If unable to remove for some reason return True."""
    if os.path.islink(filename):
        os.remove(filename)
        return False

    # Doesn't exist
    if not (os.path.isdir(filename) or os.path.isfile(filename)):
        return False

    if not overwrite:
        logger.warning(
            '%s exists and is not a symlink, Cowardly refusing to remove.',
            filename)
        return True

    if os.path.isdir(filename):
        rmtree(filename)
    else:
        os.remove(filename)

    return False


class DotfileRepo:  # pylint: disable=too-many-instance-attributes
    """
    A dotfile repo is a git repository storing dotfiles.

    This class handles all syncing and linking of a dotfile repository.
    It should not normally be used directly and instead one of Module
    or Profile should be used.
    """

    def __init__(self, where, target_dir=os.getenv('HOME')):
        self.config = None
        self.where = where
        self.target_dir = target_dir
        self.commit_msg = os.getenv('DFM_GIT_COMMIT_MSG', '')
        self.name = os.path.basename(where)

        self.files = []

        for root, dirs, files in os.walk(where):
            dirs[:] = [d for d in dirs if d != '.git']
            self.files += [os.path.join(root, f) for f in files]

        self.mappings = DEFAULT_MAPPINGS
        self.links = []
        self.hooks = {}

        dotdfm = os.path.join(where, '.dfm.yml')
        if not os.path.isfile(dotdfm):
            return

        with open(dotdfm) as dfmconfig:
            # This means it's a newer version of pyyaml
            if hasattr(yaml, "FullLoader"):
                self.config = yaml.load(dfmconfig, Loader=yaml.FullLoader)
            else:
                self.config = yaml.load(dfmconfig)

        # This indicates an empty config file
        if self.config is None:
            return

        self.target_dir = self.config.get('target_dir', self.target_dir)
        self.commit_msg = self.config.get('commit_msg', self.commit_msg)
        self.hooks = self.config.get('hooks', {})
        self.mappings = self.mappings + [
            Mapping.from_dict(mod) for mod in self.config.get('mappings', [])
        ]

    def link(self, dry_run=False, overwrite=False):
        """
        Link this profile to self.target_dir

        If the destination of a link is missing intervening
        directories this function will attempt to create them.
        """
        if not dry_run:
            self.run_hook('before_link')

        if not self.links:
            self._generate_links()

        for link in self.links:
            logger.info('Linking %s to %s', link['src'], link['dst'])
            if dry_run:
                continue

            if unable_to_remove(link['dst'], overwrite=overwrite):
                continue

            os.makedirs(os.path.dirname(link['dst']), exist_ok=True)
            os.symlink(**link)

        if not dry_run:
            self.run_hook('after_link')

        return self.links

    def _git(self, cmd, cwd=False):
        """
        Run the git subcommand 'cmd' in this dotfile repo.

        Sends all output and input to sys.stdout / sys.stdin.
        cmd should be a string and will be split using shlex.split.

        If cwd is set to None or a string then it will be passed to
        Popen constructor as the cwd argument. Otherwise the cwd for
        the process will be the location of the dotfile repo. Most
        often you will not want to set this.
        """
        try:
            if cwd or cwd is None:
                cwd = cwd
            else:
                cwd = self.where

            proc = subprocess.Popen(
                ['git'] + shlex.split(cmd),
                cwd=cwd,
                stdin=sys.stdin,
                stdout=sys.stdout,
                stderr=sys.stderr)
            proc.wait()
        except OSError as os_err:
            logger.error('problem runing git %s: %s', cmd, os_err)
            sys.exit(1)

    def _is_dirty(self):
        """
        Return the output of 'git status --porcelain'.

        This is useful because in Python an empty string is False. The
        --porcelain flag prints nothing if the git repo is not in a dirty state.
        Therefore 'if self._is_dirty()' will behave as expected.
        """
        try:
            return subprocess.check_output(
                ['git', 'status', '--porcelain'], cwd=self.where)
        # Something unexpected happened while running git so let's
        # assume we can't run anymore git commands and skip trying to
        # sync.
        except OSError:
            return False

    def run_hook(self, name):
        """Run the hook with name."""
        commands = self.hooks.get(name, [])
        for command in commands:
            try:
                subprocess.call(
                    ['/bin/sh', '-c', command],
                    cwd=self.where,
                    stdin=sys.stdin,
                    stdout=sys.stdout,
                    stderr=sys.stderr,
                )
            except subprocess.CalledProcessError as proc_err:
                logger.error('command %s exited with non-zero error: %s',
                             command, proc_err)

    def sync(self):
        """Sync this profile with git."""
        self.run_hook('before_sync')

        dirty = self._is_dirty()
        if dirty:
            if not self.commit_msg and self.config.get(
                    'prompt_for_commit_message'):
                self.commit_msg = input('Commit message: ')
            elif not self.commit_msg:
                self.commit_msg = 'Files managed by DFM! https://github.com/chasinglogic/dfm'

            self._git('add --all')
            self._git('commit -m "{}"'.format(self.commit_msg))
        self._git('pull --rebase origin master')
        if dirty:
            self._git('push origin master')

        self.run_hook('after_sync')

    def _generate_link(self, filename):
        """Dotfile-ifies a filename"""
        # Get the absolute path to src
        src = os.path.abspath(filename)
        dest = src.replace(self.where, '')

        # self.where does not always contain a trailing slash
        # This removes a leading slash from the front of dest if where
        # does not contain the trailing slash.
        if dest.startswith('/'):
            dest = dest[1:]

        dest = os.path.join(self.target_dir, dest)

        for mapping in self.mappings:
            # If the mapping doesn't match skip to the next one
            if not mapping.matches(filename):
                continue

            # If the mapping did match and is a skip mapping then end
            # function without adding a link to self.links
            if mapping.skip:
                return

            # Replace self.target_dir with the mapping target_dir
            dest = dest.replace(self.target_dir, mapping.target_dir)

        self.links.append({
            'src': src,
            'dst': dest,
            'target_is_directory': os.path.isdir(src)
        })

    def _generate_links(self):
        """
        Generate a list of kwargs for os.link.

        All required arguments for os.link will always be provided and
        optional arguments as required.
        """
        for dotfile in self.files:
            self._generate_link(dotfile)


class Module(DotfileRepo):
    """
    Module is a DotfileRepo that has additional options for syncing and linking.

    Module provides a new option for syncing called 'pull_only' which
    will not push to the remote repo.

    Additionally a Module has a known location on the filesystem,
    either auto-generated or manually specified, and if not found will
    attempt to clone the repository provided as the argument 'repo'
    into that location.

    Module also feeds the pre or post property up to it's parent
    profile to determine when it should be linked in relation to that
    profile.
    """

    def __init__(self, *args, **kwargs):
        self.repo = kwargs.pop('repo')
        self.name = kwargs.pop('name', '')
        if not self.name:
            self.name = self.repo.split('/')[-1]
        self.pull_only = kwargs.pop('pull_only', False)
        self.link_mode = kwargs.pop('link', 'post')
        self.location = kwargs.pop('location', '')
        self.location = self.location.replace('~', os.getenv('HOME'))
        if not self.location:
            module_dir = os.path.join(dfm_dir(), 'modules')
            if not os.path.isdir(module_dir):
                os.makedirs(module_dir)
            self.location = os.path.join(module_dir, self.name)

        if not os.path.isdir(
                self.location) and not os.getenv('DFM_DISABLE_MODULES'):
            self._git('clone {} {}'.format(self.repo, self.location), cwd=None)

        kwargs['where'] = self.location
        super().__init__(*args, **kwargs)

    def sync(self):
        """Sync this repo using git, if self.pull_only will only pull updates."""
        if self.pull_only:
            self._git('pull --rebase origin master')
            return

        super().sync()

    def link(self, dry_run=False, overwrite=False):
        """Wrap super()._generate_links"""
        if self.link_mode == 'none':
            return []

        return super().link(dry_run=dry_run, overwrite=overwrite)

    @property
    def pre(self):
        """If True this module should be linked before the parent Profile."""
        return self.link_mode == 'pre'

    @property
    def post(self):
        """
        If True this module should be linked after the parent Profile.

        This is useful for when you want files from a module to
        overwrite those from it's parent Profile.
        """
        return self.link_mode == 'post'

    @classmethod
    def from_dict(cls, config):
        """Return a Module from the config dictionary"""
        return cls(**config)


class Profile(DotfileRepo):
    """Profile is a DotfileRepo that supports modules."""

    def __init__(self,
                 where,
                 always_sync_modules=False,
                 target_dir=os.getenv('HOME')):

        super().__init__(where, target_dir=target_dir)
        self.always_sync_modules = always_sync_modules
        self.modules = []

        if self.config is None:
            return

        self.always_sync_modules = self.config.get('always_sync_modules',
                                                   self.always_sync_modules)
        self.modules = [
            Module.from_dict(mod) for mod in self.config.get('modules', [])
        ]

    def sync(self, skip_modules=False):  # pylint: disable=arguments-differ
        """
        Sync this profile and all modules.

        If skip_modules is True modules will not be synced.
        """
        print('{}:'.format(self.where))
        super().sync()

        if skip_modules:
            return

        for module in self.modules:
            print('\n{}:'.format(module.where))
            module.sync()

    def _generate_links(self):
        """Add module support to DotfileRepo's _generate_links."""
        for module in self.modules:
            if module.pre:
                self.links += module.link(dry_run=True)

        super()._generate_links()

        for module in self.modules:
            if module.post:
                self.links += module.link(dry_run=True)
