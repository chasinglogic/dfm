import re
from os.path import expanduser


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

    @classmethod
    def from_dict(cls, config):
        """Return a Mapping from the config dictionary"""
        return cls(**config)

    def matches(self, path):
        """Determine if this mapping matches path."""
        return self.rgx.search(path)


DEFAULT_MAPPINGS = [
    Mapping(r"\/\.git\/", skip=True,),
    Mapping(r"\/.gitignore$", skip=True,),
    Mapping(r"\/.ggitignore$", dest=".gitignore"),
    Mapping(r"\/LICENSE(\.md)?$", skip=True,),
    Mapping(r"\/\.dfm\.yml$", skip=True,),
    Mapping(r"\/README(\.md)?$", skip=True,),
]
