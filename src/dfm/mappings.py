"""Map files to non-normal locations."""

import os
import platform
import re
from os.path import expanduser

CUR_OS = platform.system()


class Mapping:
    """
    Maps a filename to a new destination.

    Allows for dotfiles to be skipped, or redirected to a target
    directory other than '$HOME'.
    """

    def __init__(self, match, dest="", target_dir="", skip=False, target_os=None):
        self.match = match
        self.target_dir = expanduser(target_dir)
        self.dest = expanduser(dest)
        if target_os is not None:
            self.target_os = target_os
        else:
            self.target_os = []
        self.skip = skip
        self.rgx = re.compile(match)

    def on_target_os(self):
        """Return a boolean indicating if we're on the target OS for this Mapping."""
        return (
            isinstance(self.target_os, list) and CUR_OS in self.target_os
        ) or self.target_os == CUR_OS

    def should_skip(self):
        """Determine if the file matching this Mapping should be skipped."""
        if not self.skip:
            return False

        # We aren't an OS-specific mapping and are a skip mapping so return should skip.
        if not self.target_os:
            return True

        # We are a skip mapping for this OS so return True
        if self.on_target_os():
            return True

        # We are a skip mapping but we aren't on the target OS return False so the file
        # is not skipped.
        return False

    def replace(self, dest, target_dir):
        """Return the new destination for link based on this Mapping."""
        if self.target_os and not self.on_target_os():
            return dest

        if self.dest:
            new_dest = expanduser(self.dest)
            if new_dest[0] == os.path.pathsep:
                return new_dest

            return os.path.join(target_dir, new_dest)

        if self.target_dir:
            return dest.replace(target_dir, self.target_dir)

        return dest

    @classmethod
    def from_dict(cls, config):
        """Return a Mapping from the config dictionary."""
        return cls(**config)

    @classmethod
    def from_config(cls, config):
        """Load Mappings from config."""
        mappings = [cls(**mapping) for mapping in config.pop("mappings", [])]
        mappings.extend(DEFAULT_MAPPINGS)
        return mappings

    def matches(self, path):
        """Determine if this mapping matches path."""
        return self.rgx.search(path) is not None

    def __repr__(self):
        if self.dest:
            new_dest = self.dest
        elif self.target_dir:
            new_dest = self.target_dir + os.path.pathsep
        elif self.skip:
            new_dest = "SKIP"
        else:
            new_dest = "UNKNOWN"

        return "Mapping({from_match} -> {to}, os={os})".format(
            from_match=self.match,
            to=new_dest,
            os=self.target_os,
        )


DEFAULT_MAPPINGS = [
    Mapping(
        r"\/\.git\/",
        skip=True,
    ),
    Mapping(
        r"\/.gitignore$",
        skip=True,
    ),
    Mapping(r"\/.ggitignore$", dest=".gitignore"),
    Mapping(
        r"\/LICENSE(\.md)?$",
        skip=True,
    ),
    Mapping(
        r"\/\.dfm\.yml$",
        skip=True,
    ),
    Mapping(
        r"\/README(\.md)?$",
        skip=True,
    ),
]
