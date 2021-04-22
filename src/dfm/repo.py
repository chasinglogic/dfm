"""Manage git commands on a repository."""

import logging
import os
import shlex
import subprocess
import sys

logger = logging.getLogger(__name__)


class DotfileRepo:  # pylint: disable=too-many-instance-attributes
    """
    A dotfile repo is a git repository storing dotfiles.

    This class handles all syncing and git operations of a dotfile repository. It should
    not normally be used directly and instead Profile should be used.
    """

    def __init__(
        self,
        where,
        commit_msg,
        prompt_for_commit_message=False,
    ):
        self.where = where
        self.commit_msg = commit_msg
        self.prompt_for_commit_message = prompt_for_commit_message

    @classmethod
    def from_config(cls, where, config):
        """Load DotfileRepo settings from config."""
        return cls(
            where=where,
            commit_msg=config.pop("commit_msg", os.getenv("DFMgit_COMMIT_MSG")),
            prompt_for_commit_message=config.pop("prompt_for_commit_message", False),
        )

    def init(self):
        """Initialize the git repository."""
        self.git("init")

    def get_remote(self):  # pylint: disable=no-self-use
        """Return the remote url for origin."""
        return subprocess.check_output(
            ["git", "remote", "get-url", "origin"],
            cwd=self.where,
        ).decode("utf-8")

    def git(self, cmd, cwd=False, dry_run=False):
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
            if cwd is False:
                cwd = self.where

            args = ["git"] + shlex.split(cmd)

            if dry_run:
                logger.info('Running: "%s" in %s', " ".join(args), cwd)
                return

            proc = subprocess.Popen(
                args,
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
                ["git", "status", "--porcelain"],
                cwd=self.where,
            )
        # Something unexpected happened while running git so let's
        # assume we can't run anymore git commands and skip trying to
        # sync.
        except OSError:
            return False

    def sync(self, commit_msg="", dry_run=False):
        """Sync this profile with git."""
        logger.info("Syncing: %s", self.where)
        dirty = self._is_dirty()
        if dirty:
            self.git("--no-pager diff", dry_run=dry_run)

            if not self.commit_msg and dry_run:
                self.commit_msg = "noop"

            elif not self.commit_msg and commit_msg:
                self.commit_msg = commit_msg

            elif not self.commit_msg and self.prompt_for_commit_message:
                self.commit_msg = input("Commit message: ")

            elif not self.commit_msg:
                self.commit_msg = (
                    "Files managed by DFM! https://github.com/chasinglogic/dfm"
                )

            self.git("add --all", dry_run=dry_run)
            self.git(
                'commit -m "{}"'.format(self.commit_msg),
                dry_run=dry_run,
            )

        if self._has_origin():
            self.git(
                "pull --rebase origin master",
                dry_run=dry_run,
            )
            if dirty:
                self.git(
                    "push origin master",
                    dry_run=dry_run,
                )

    def _has_origin(self):
        """Return a boolean indicating if a remote named origin is available."""
        try:
            output = subprocess.check_output(["git", "remote", "-v"], cwd=self.where)
            return b"origin" in output
        except OSError:
            return False
