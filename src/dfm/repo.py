"""Manage git commands on a repository."""

import logging
import os
import shlex
import subprocess
import sys
import git
from git.exc import InvalidGitRepositoryError, NoSuchPathError, GitCommandError, BadName

logger = logging.getLogger(__name__)


class DotfileRepo:  # pylint: disable=too-many-instance-attributes
    """
    A dotfile repo is a git repository storing dotfiles.

    This class handles all syncing and git operations of a dotfile repository. It should
    not normally be used directly and instead Profile should be used.
    """

    def __init__(
        self,
        path,
        remote_url=None,
        branch_name=None,
        commit_msg=None,
        prompt_for_commit_message=False,
    ):
        self.local_path = path
        self.commit_msg = commit_msg
        self.prompt_for_commit_message = prompt_for_commit_message
        #
        self._git_repo = None

        self._init_branch_name = branch_name
        self._remote_url = remote_url

        try:
            self._git_repo = git.Repo(self.local_path)
            self._try_set_remote()
            self._try_set_branch()
        except (InvalidGitRepositoryError,NoSuchPathError):
            # repo could not initialised (path does not exist, for instance)
            #, will have to call initialise()
            self._git_repo = None

    def _try_set_remote(self):
        self._remote_name = self._get_remote_name(remote_url=self._remote_url)
        # create the remote if we have a url but no name
        if self._remote_url is not None and self._remote_name is None:
            self._remote_name = 'origin'
        assert self._git_repo is not None
        if self._remote_url is not None and self.get_remote() is None:
            remote = self._git_repo.create_remote(self._remote_name, self._remote_url)
            assert remote.exists()
        return self.get_remote()

    def _try_set_branch(self):
        self._branch_name = self._get_branch_name(remote_name=self._remote_name, branch_name=self._init_branch_name)
        assert self._branch_name is not None
        remote = self.get_remote()
        branch = self.get_branch()
        if branch is None:
            try:
                if remote is not None:
                    remote_branch = remote.refs[self._branch_name]
                    branch = self._git_repo.create_head(self._branch_name, remote_branch)
                else:
                    branch = self._git_repo.create_head(self._branch_name)
            except (GitCommandError, BadName):
                # TODO: BadName may be a bug??
                branch = None
        if branch is not None:
            if remote is not None:
                remote_branch = remote.refs[self._branch_name]
                branch.set_tracking_branch(remote_branch)
        return self.get_branch()



    def is_initialised(self):
        return self._git_repo is not None

    def initialise(self):
        """
            Creates git repo
        """
        self._git_repo = git.Repo.init(self.local_path)
        remote = self._try_set_remote()
        if remote is not None:
            remote.fetch()
        branch = self._try_set_branch()
        if branch is not None:
            branch.checkout()

    def get_remote(self):
        if not self.is_initialised():
            return None
        has_remote = False
        try:
            remote = self._git_repo.remote(self._remote_name)
            has_remote = remote.exists()
        except ValueError:
            has_remote = False
        return remote if has_remote else None

    def get_remote_url(self):
        remote = self.get_remote()
        if remote is not None:
            url = next(remote.urls)
            assert url is not None
            return url
        else:
            return self._remote_url

    def get_branch(self):
        if not self.is_initialised():
            return None
        branch = None
        if self._branch_name is not None:
            try:
                branch = self._git_repo.heads[self._branch_name]
            except IndexError:
                pass
        return branch

    def _get_remote_name(self, remote_url=None):
        # take the first remote by default
        # we get the name of the remote that we are interest about
        has_remote = len(self._git_repo.remotes) > 0
        remote_name = None
        if remote_url is None:
            # if no remote was explicitly given,
            # we take the first one we found if it exists.
            if has_remote:
                remote_name = self._git_repo.remotes[0].name
        else:
            # find the name of the remote with the url
            for remote in self._git_repo.remotes:
                if remote.url == remote_url:
                    remote_name = remote.name
        return remote_name

    def _get_branch_name(self, remote_name=None, branch_name=None):
        active_branch_name = None
        remote = self.get_remote()
        has_remote = remote is not None
        has_active_branch = False
        try:
            active_branch_name = self._git_repo.active_branch.name
            has_active_branch = True
        except TypeError:
            has_active_branch = False
        has_branch = len(self._git_repo.heads)
        if branch_name is None:
            # take the active branch name if it exists
            if has_active_branch:
                branch_name = active_branch_name
            elif has_branch:
                # take the first branch that exists
                branch_name = self._git_repo.heads[0].name
            elif has_remote and len(remote.repo.heads) > 0:
                branch_name = remote.repo.heads[0].name
            else:
                # use a default one.
                branch_name = 'main'
        return branch_name

    @classmethod
    def from_config(cls, path, config: dict):
        """Load DotfileRepo settings from config."""
        return cls(
            path=path,
            remote_url=config.get("repo"),
            branch_name=config.get("branch"),
            commit_msg=config.pop("commit_msg", os.getenv("DFMgit_COMMIT_MSG")),
            prompt_for_commit_message=config.pop("prompt_for_commit_message", False),
        )


    def fetch(self, **kwargs):
        # TODO: progress info
        remote = self.get_remote()
        assert remote is not None
        remote.fetch(**kwargs)

    def pull(self, **kwargs):
        # TODO: progress info
        remote = self.get_remote()
        assert remote is not None
        remote.pull(**kwargs)

    def push(self, **kwargs):
        # TODO: progress info
        remote = self.get_remote()
        assert remote is not None
        remote.push(**kwargs)

    def query_commit_message(self, default_msg=None, dry_run=False):
        commit_msg = self.commit_msg
        if not commit_msg:
            if dry_run:
                commit_msg = "noop"
            elif default_msg:
                commit_msg = default_msg
            elif self.prompt_for_commit_message:
                commit_msg = input("Commit message: ")
            else:
                commit_msg = "Files managed by DFM! https://github.com/chasinglogic/dfm"
        return commit_msg


    def sync(self, commit_msg="", dry_run=False):
        """Sync this profile with git."""
        logger.info("Syncing: %s", self.local_path)
        remote = self._try_set_remote()
        if remote is not None:
            self.fetch()
        branch = self._try_set_branch()
        if branch is not None:
            branch.checkout()
        if self._git_repo.is_dirty(untracked_files=True):
            # repo_index.diff(
            commit_msg = self.query_commit_message(commit_msg, dry_run=dry_run)
            repo_index = self._git_repo.index
            if not dry_run:
                repo_index.add('./')
                commit = repo_index.commit(commit_msg)
                # we try again to create a branch, now that we have a commit
                branch = self._try_set_branch()
                assert branch is not None
        if remote is not None:
            self.pull(rebase=True, dry_run=dry_run)
            self.push(dry_run=dry_run)
