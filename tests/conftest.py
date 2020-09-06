"""Common test fixtures"""

import os

import pytest


@pytest.fixture
def dotfile_dir(tmpdir):
    """Return a pre-populated dotfile directory with some files."""

    DEFAULT_DOTFILES = [
        ".vimrc",
        ".bashrc",
        ".emacs",
        ".gitignore",
        ".ggitignore",
        ".emacs.d/init.el",
        ".git/HEAD",
    ]

    def touch(dotfile):
        open(dotfile, "w").close()

    def create_dotfile_dir(dotfiles=None):
        if dotfiles is None:
            dotfiles = DEFAULT_DOTFILES

        for dotfile in dotfiles:
            elements = os.path.split(dotfile)
            directories, file = elements[:-1], elements[-1]
            if directories:
                df_dir = os.path.join(tmpdir, *directories)
                os.makedirs(df_dir, exist_ok=True)
            else:
                df_dir = tmpdir

            touch(os.path.join(df_dir, file))

        return (
            [os.path.join(tmpdir, df) for df in dotfiles if not df.startswith(".git/")],
            tmpdir,
        )

    return create_dotfile_dir


@pytest.fixture
def dotdfm(dotfile_dir):
    """Create a dotdfm with passed content"""

    def create_dotdfm(content="", dotfiles=None):
        dotfiles, directory = dotfile_dir(dotfiles=dotfiles)

        with open(os.path.join(directory, ".dfm.yml"), "w") as cfg:
            cfg.write(content)

        return dotfiles, directory

    return create_dotdfm
