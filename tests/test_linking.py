"""Test linking"""

import os
from operator import itemgetter
from tempfile import TemporaryDirectory

from dfm.dotfile import Profile, xdg_dir


def setup_module():
    """Set DFM_CONFIG_DIR for this module's tests."""
    os.environ["DFM_CONFIG_DIR"] = TemporaryDirectory().name


def test_translation(dotdfm):
    """Test that the generated links properly translate names."""
    _, directory = dotdfm()
    profile = Profile(str(directory))
    links = profile.link(dry_run=True)
    expected_links = [
        {
            "src": os.path.join(directory, ".vimrc"),
            "dst": os.path.join(os.getenv("HOME"), ".vimrc"),
        },
        {
            "src": os.path.join(directory, ".bashrc"),
            "dst": os.path.join(os.getenv("HOME"), ".bashrc"),
        },
        {
            "src": os.path.join(directory, ".emacs"),
            "dst": os.path.join(os.getenv("HOME"), ".emacs"),
        },
        {
            "src": os.path.join(directory, ".ggitignore"),
            "dst": os.path.join(os.getenv("HOME"), ".gitignore"),
        },
        {
            "src": os.path.join(directory, ".emacs.d", "dotfile"),
            "dst": os.path.join(os.getenv("HOME"), ".emacs.d", "dotfile"),
        },
    ]

    for (link, expected) in zip(
        sorted(links, key=itemgetter("src")),
        sorted(expected_links, key=itemgetter("src")),
    ):
        assert link["src"] == expected["src"]
        assert link["dst"] == expected["dst"]


def test_linking(dotdfm):
    """Test that the profile properly creates the links."""
    _, directory = dotdfm()
    with TemporaryDirectory() as target:
        profile = Profile(str(directory), target_dir=target)
        links = profile.link()
        # Use a set comprehension since the target would not contain duplicates
        # since they would be overwritten with the last occuring link.
        dest = sorted(list({os.path.basename(x["dst"]) for x in links}))
        created = []
        for root, dirs, files in os.walk(target):
            dirs[:] = [d for d in dirs if d != ".git"]
            created += [f for f in files if os.path.islink(os.path.join(root, f))]

        created = sorted(created)
        assert dest == created


def test_xdg_config(tmpdir):
    dotfiles = tmpdir.mkdir("dotfiles")
    initvim = dotfiles.mkdir(".config").mkdir("nvim").join("init.vim")
    open(initvim, "w").close()
    profile = Profile(str(dotfiles))
    links = profile.link(dry_run=True)
    assert links == [
        {
            "src": os.path.join(dotfiles, ".config", "nvim", "init.vim"),
            "dst": os.path.join(xdg_dir(), "nvim", "init.vim"),
        }
    ]
