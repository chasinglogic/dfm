import re
import platform
import os
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
        return (
            isinstance(self.target_os, list) and CUR_OS in self.target_os
        ) or self.target_os == CUR_OS

    def should_skip(self):
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
        """Return a Mapping from the config dictionary"""
        return cls(**config)

    def matches(self, path):
        """Determine if this mapping matches path."""
        return self.rgx.search(path) is not None

    def __repr__(self):
        if self.dest:
            to = self.dest
        elif self.target_dir:
            to = self.target_dir + os.path.pathsep
        elif self.skip:
            to = "SKIP"
        else:
            to = "UNKNOWN"

        return "Mapping({from_match} -> {to}, os={os})".format(
            from_match=self.match,
            to=to,
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
