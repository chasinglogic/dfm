"""Hook management for Profiles."""

import logging
import shlex
import subprocess
import sys

logger = logging.getLogger(__name__)


class Hooks:
    """Manages and executes hooks."""

    def __init__(self, where, hooks=None):
        self.where = where
        self.hooks = hooks if hooks is not None else {}

    @classmethod
    def from_config(cls, where, config):
        """Load Hooks from a .dfm.yml config file."""
        return cls(where=where, hooks=config.pop("hooks", {}))

    def run_hook(self, name, dry_run=False):
        """
        Run the hook with name.

        If dry_run is given the command will not be actually executed.
        """
        commands = self.hooks.get(name, [])
        for command in commands:
            if isinstance(command, dict):
                interpreter = shlex.split(command.get("interpreter", "/bin/sh -c"))
                script = command.get("script", "")
            else:
                interpreter = ["/bin/sh", "-c"]
                script = command

            if not script:
                logger.warning(
                    "Found an empty script for hook: %s, skipping",
                    name,
                )
                continue

            if dry_run:
                logger.debug("Running hook %s command: %s", name, command)
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
                    "command %s exited with non-zero error: %s",
                    command,
                    proc_err,
                )
