"""Dotfile profile and module management."""

import logging
import os
import warnings
from subprocess import call

import yaml

from dfm.config import dfm_dir
from dfm.hooks import Hooks
from dfm.links import LinkManager
from dfm.repo import DotfileRepo

logger = logging.getLogger(__name__)


DEPRECATED_CONFIG_OPTIONS = [
    "always_sync_modules",
]


def get_name(url):
    """
    Generate a profile name based on the git url.

    This directly corresponds to the second to last element in the URL.  For
    example: https://github.com/chasinglogic/dotfiles would be 'chasinglogic'

    In the case of an ssh url or other url will correspond to the
    first path element. For example:
    git@github.com:chasinglogic/dotfiles would be 'chasinglogic'
    """
    try:
        if not url:
            return ""

        if "@" in url:
            return url.split(":")[-1].split("/")[0]

        return url.split("/")[-2]
    # Any kind of exception i.e. IndexError or failure to split we
    # just return nothing.
    except IndexError:
        return ""


class Profile:  # pylint: disable=too-many-instance-attributes
    """
    Profile is a dotfile profile.

    It is the aggregation of a single dotfile repo and all of it's modules.

    A Module (which is a Profile instance) has a known location on the
    filesystem, either auto-generated or manually specified, and if
    not found will attempt to clone the repository provided as the
    argument 'repo' into that location.

    Modules provide a new option for syncing called 'pull_only' which
    will not push to the remote repo. This is useful for third-party dotfiles like
    oh-my-zsh.

    A Module also feeds the pre or post property up to it's parent
    Profile to determine when it should be linked in relation to that
    Profile.
    """

    def __init__(  # pylint: disable=too-many-arguments
        self,
        df_repo,
        link_manager,
        hooks,
        repo="",
        repository="",
        name="",
        pull_only=False,
        link="post",
        branch="master",
        modules=None,
    ):
        self.df_repo = df_repo
        self.link_manager = link_manager
        self.hooks = hooks
        self.pull_only = pull_only
        self.link_mode = link
        self.repo = repo if repo else repository
        if not self.repo:
            self.repo = self.df_repo.get_remote()
        self.modules = modules if modules is not None else []
        self.branch = branch
        self.name = name or get_name(self.repo)

    def sync(self, commit_msg="", dry_run=False, skip_modules=False):
        """
        Sync this profile and all modules using git.

        If self.pull_only will only pull updates.

        If skip_modules is True modules will not be synced.
        """
        self.hooks.run_hook("before_sync", dry_run=dry_run)
        if self.pull_only:
            self.df_repo.git(
                "pull --rebase origin {}".format(self.branch),
                dry_run=dry_run,
            )
        else:
            self.df_repo.sync(
                dry_run=dry_run,
                commit_msg=commit_msg,
            )
        self.hooks.run_hook("after_sync", dry_run=dry_run)

        if skip_modules:
            return

        for module in self.modules:
            module.sync(dry_run=dry_run)

    def link(self, dry_run=False, overwrite=False):
        """Link the dotfiles link for this profile and its modules."""
        for module in self.pre_link_modules:
            module.link(dry_run=dry_run, overwrite=overwrite)

        if self.should_link:
            self.hooks.run_hook("before_link", dry_run=dry_run)
            self.link_manager.link(dry_run, overwrite=overwrite)
            self.hooks.run_hook("after_link", dry_run=dry_run)

        for module in self.post_link_modules:
            module.link(dry_run=dry_run, overwrite=overwrite)

    @property
    def pre_link_modules(self):
        """Return an iterator for all pre link modules."""
        for module in self.modules:
            if module.link_pre:
                yield module

    @property
    def post_link_modules(self):
        """Return an iterator for all post link modules."""
        for module in self.modules:
            if module.link_post:
                yield module

    @property
    def link_pre(self):
        """If True this module should be linked before the parent Profile."""
        return self.link_mode == "pre"

    @property
    def should_link(self):
        """If True this profile should be linked."""
        return self.link_mode != "none"

    @property
    def link_post(self):
        """
        If True this module should be linked after the parent Profile.

        This is useful for when you want files from a module to
        overwrite those from it's parent Profile.
        """
        return self.link_mode == "post"

    @classmethod
    def default(cls, where, extras=None):
        """
        Return a default Profile for where.

        This is used when a config is not present or empty.
        """
        if extras is None:
            extras = {}
        df_repo = DotfileRepo.from_config(where, extras)
        link_manager = LinkManager.from_config(where, extras)
        hooks = Hooks.from_config(where, extras)
        return cls(df_repo, link_manager, hooks, **extras)

    @classmethod
    def load(cls, where, extras=None):
        """
        Load a profile from where.

        This loads the .dfm.yml file and configures the Profile and any modules
        accordingly.
        """
        dotdfm = os.path.join(where, ".dfm.yml")
        if not os.path.isfile(dotdfm):
            return cls.default(where, extras=extras)

        with open(dotdfm) as dfmconfig:
            # This means it's a newer version of pyyaml
            if hasattr(yaml, "FullLoader"):
                config = yaml.load(dfmconfig, Loader=yaml.FullLoader)
            else:
                config = yaml.load(dfmconfig)

        # This indicates an empty config file
        if config is None:
            return cls.default(where, extras=extras)

        if extras:
            config.update(extras)

        for key in DEPRECATED_CONFIG_OPTIONS:
            value = config.pop(key, None)
            if value is not None:
                warnings.warn(
                    "The config option {} has been deprecated, ignoring.".format(
                        key,
                    ),
                )

        modules = [cls.load_module(mod) for mod in config.pop("modules", [])]
        df_repo = DotfileRepo.from_config(where, config)
        link_manager = LinkManager.from_config(where, config)
        hooks = Hooks.from_config(where, config)
        profile = cls(
            df_repo=df_repo,
            link_manager=link_manager,
            hooks=hooks,
            modules=modules,
            **config,
        )
        logger.debug(
            "Loaded %s: %s",
            "profile" if not extras else "module",
            profile.name,
        )
        return profile

    @classmethod
    def load_module(cls, config):
        """
        Load a module as a Profile.

        This handles the necessary extra setup and config manipulation require to
        properly initialise a module.
        """
        if "repo" not in config:
            config["repo"] = config.pop("repository")
        name = config.get("name", get_name(config["repo"]))
        location = os.path.expanduser(config.pop("location", ""))
        if not location:
            module_dir = os.path.join(dfm_dir(), "modules")
            if not os.path.isdir(module_dir):
                os.makedirs(module_dir)
            location = os.path.join(module_dir, name)

        if not os.path.isdir(location) and not os.getenv("DFM_DISABLE_MODULES"):
            call(
                [
                    "git",
                    "clone",
                    "--single-branch",
                    "--branch",
                    config.get("branch", "master"),
                    config["repo"],
                    location,
                ],
            )

        return cls.load(location, extras=config)

    @classmethod
    def new(cls, where):
        """Create and initialise a new Profile at where."""
        profile = cls.default(where)
        profile.df_repo.init()
        return profile

    @classmethod
    def from_dict(cls, config):
        """Return a Module from the config dictionary"""
        return cls(**config)
