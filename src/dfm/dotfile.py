"""Translate file names to the appropriate targets."""

import logging
import os
import shlex
import subprocess
import sys
import platform
from shutil import rmtree

import yaml

from dfm.mappings import Mapping, DEFAULT_MAPPINGS

logger = logging.getLogger(__name__)

CUR_OS = platform.system()


def xdg_dir():
    """Return the XDG_CONFIG_HOME or default."""
    if os.getenv("XDG_CONFIG_HOME"):
        return os.getenv("XDG_CONFIG_HOME")
    return os.path.join(os.getenv("HOME"), ".config")


def dfm_dir():
    """Return the dfm configuration / state directory."""
    if os.getenv("DFM_CONFIG_DIR"):
        return os.getenv("DFM_CONFIG_DIR")
    return os.path.join(xdg_dir(), "dfm")


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
            "%s exists and is not a symlink, Cowardly refusing to remove.", filename
        )
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
    It should not normally be used directly and instead Profile should
    be used.
    """

    def __init__(
        self,
        where,
        target_dir=os.getenv("HOME"),
        commit_msg=os.getenv("DFM_GIT_COMMIT_MSG", ""),
    ):
        self.config = None
        self.where = where
        self.target_dir = target_dir
        self.commit_msg = commit_msg
        self.name = os.path.basename(where)

        self.files = []

        self.mappings = DEFAULT_MAPPINGS
        self.links = []
        self.hooks = {}

        dotdfm = os.path.join(where, ".dfm.yml")
        if not os.path.isfile(dotdfm):
            self.config = {}
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

        self.target_dir = self.config.get("target_dir", self.target_dir)
        self.commit_msg = self.config.get("commit_msg", self.commit_msg)
        self.hooks = self.config.pop("hooks", {})
        self.mappings = self.mappings + [
            Mapping.from_dict(mod) for mod in self.config.pop("mappings", [])
        ]

    def init(self):
        """Initialize the git repository."""
        self._git("init")

    def link(self, dry_run=False, overwrite=False):
        """
        Link this profile to self.target_dir

        If the destination of a link is missing intervening
        directories this function will attempt to create them.
        """
        if not dry_run:
            self.run_hook("before_link")

        if not self.links:
            self._generate_links()

        for link in self.links:
            logger.info("Linking %s to %s", link["src"], link["dst"])
            if dry_run:
                continue

            if unable_to_remove(link["dst"], overwrite=overwrite):
                continue

            os.makedirs(os.path.dirname(link["dst"]), exist_ok=True)
            os.symlink(**link)

        if not dry_run:
            self.run_hook("after_link")

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
            if not cwd:
                cwd = self.where

            proc = subprocess.Popen(
                ["git"] + shlex.split(cmd),
                cwd=cwd,
                stdin=sys.stdin,
                stdout=sys.stdout,
                stderr=sys.stderr,
            )
            proc.wait()
        except OSError as os_err:
            logger.error("problem runing git %s: %s", cmd, os_err)
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
                ["git", "status", "--porcelain"], cwd=self.where
            )
        # Something unexpected happened while running git so let's
        # assume we can't run anymore git commands and skip trying to
        # sync.
        except OSError:
            return False

    def run_hook(self, name):
        """Run the hook with name."""
        commands = self.hooks.get(name, [])
        for command in commands:
            if isinstance(command, dict):
                interpreter = shlex.split(command.get("interpreter", "/bin/sh -c"))
                script = command.get("script", "")
            else:
                interpreter = ["/bin/sh", "-c"]
                script = command

            if not script:
                print("Found an empty script for hook: {}, skipping".format(name))
                continue

            try:
                subprocess.call(
                    interpreter + [script],
                    cwd=self.where,
                    stdin=sys.stdin,
                    stdout=sys.stdout,
                    stderr=sys.stderr,
                )
            except subprocess.CalledProcessError as proc_err:
                logger.error(
                    "command %s exited with non-zero error: %s", command, proc_err
                )

    def sync(self):
        """Sync this profile with git."""
        self.run_hook("before_sync")

        dirty = self._is_dirty()
        if dirty:
            self._git("--no-pager diff")

            if not self.commit_msg and self.config.get("prompt_for_commit_message"):
                self.commit_msg = input("Commit message: ")
            elif not self.commit_msg:
                self.commit_msg = (
                    "Files managed by DFM! https://github.com/chasinglogic/dfm"
                )

            self._git("add --all")
            self._git('commit -m "{}"'.format(self.commit_msg))

        if self._has_origin():
            self._git("pull --rebase origin master")
            if dirty:
                self._git("push origin master")

        self.run_hook("after_sync")

    def _has_origin(self):
        """Return a boolean indicating if a remote named origin is available."""
        try:
            output = subprocess.check_output(["git", "remote", "-v"], cwd=self.where)
            return b"origin" in output
        except OSError:
            return False

    def _generate_link(self, filename):
        """Dotfile-ifies a filename"""
        # Get the absolute path to src
        src = os.path.abspath(filename)
        dest = src.replace(self.where, "")

        # self.where does not always contain a trailing slash
        # This removes a leading slash from the front of dest if where
        # does not contain the trailing slash.
        if dest.startswith("/"):
            dest = dest[1:]

        dest = os.path.join(self.target_dir, dest)

        for mapping in self.mappings:
            # If the mapping doesn't match skip to the next one
            if not mapping.matches(filename):
                continue

            # If the mapping did match and is a skip mapping then end
            # function without adding a link to self.links
            if mapping.should_skip():
                return

            dest = mapping.replace(dest, self.target_dir)

        self.links.append({"src": src, "dst": dest})

    def _find_files(self):
        """Load the files in this dotfile repository."""
        for root, dirs, files in os.walk(self.where):
            dirs[:] = [d for d in dirs if d != ".git"]
            self.files += [os.path.join(root, f) for f in files]

    def _generate_links(self):
        """
        Generate a list of kwargs for os.link.

        All required arguments for os.link will always be provided and
        optional arguments as required.
        """
        if not self.files:
            self._find_files()

        for dotfile in self.files:
            self._generate_link(dotfile)
