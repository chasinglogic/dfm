"""Common test fixtures"""

import os

import pytest


@pytest.fixture
def dotfile_dir(tmpdir):
    """Return a pre-populated dotfile directory with some files."""
    dotfiles = [
        ".vimrc",
        ".bashrc",
        ".emacs",
        ".gitignore",
        ".ggitignore",
    ]

    dotfile_dirs = [
        ".emacs.d",
        ".git",
    ]

    for dotfile in dotfiles:
        open(os.path.join(tmpdir, dotfile), "w").close()

    for dotfile in dotfile_dirs:
        df_dir = os.path.join(tmpdir, dotfile)
        os.mkdir(df_dir)
        if dotfile != ".git":
            df_dir_df = os.path.join(df_dir, "dotfile")
            dotfiles.append(df_dir_df)
            open(df_dir_df, "w").close()

    return ([os.path.join(tmpdir, df) for df in dotfiles], tmpdir)


@pytest.fixture
def dotdfm(dotfile_dir):
    """Create a dotdfm with passed content"""

    def create_dotdfm(content=""):
        dotfiles, directory = dotfile_dir

        with open(os.path.join(directory, ".dfm.yml"), "w") as cfg:
            cfg.write(content)

        return dotfiles, directory

    return create_dotdfm
