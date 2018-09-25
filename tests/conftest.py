"""Common test fixtures"""

import os

import pytest


@pytest.fixture
def dotfile_dir(tmpdir):
    dotfiles = [
        'vimrc',
        'bashrc',
        'emacs',
        'gitignore',
        '.ggitignore',
        '.gitignore',
    ]

    dotfile_dirs = [
        'emacs.d',
        '.git',
    ]

    for dotfile in dotfiles:
        open(os.path.join(tmpdir, dotfile), 'w').close()

    for dotfile in dotfile_dirs:
        os.mkdir(os.path.join(tmpdir, dotfile))

    return (dotfiles + dotfile_dirs, tmpdir)


@pytest.fixture
def dotdfm(dotfile_dir):
    """Create a dotdfm with passed content"""

    def create_dotdfm(content=''):
        dotfiles, directory = dotfile_dir

        with open(os.path.join(directory, '.dfm.yml'), 'w') as cfg:
            cfg.write(content)

        return dotfiles, directory

    return create_dotdfm
