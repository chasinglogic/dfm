"""Test linking"""

import os
import platform
from operator import itemgetter
from tempfile import TemporaryDirectory

from dfm.config import xdg_dir
from dfm.links import LinkManager
from dfm.profile import Profile


def setup_module():
    """Set DFM_CONFIG_DIR for this module's tests."""
    os.environ["DFM_CONFIG_DIR"] = TemporaryDirectory().name


def test_translation(dotdfm):
    """Test that the generated links properly translate names."""
    _, directory = dotdfm(
        """
---
mappings:
    - match: .skip_on_os
      target_os:
         - {this_os}
      skip: true
    - match: .map_to_os_name
      dest: .{this_os}
      target_os: {this_os}
    - match: .skip_on_another_os
      skip: True
      target_os: {another_os}
""".format(
            this_os=platform.system(),
            another_os="Windows" if platform.system() != "Windows" else "Linux",
        ),
        dotfiles=[
            ".vimrc",
            ".bashrc",
            ".emacs",
            ".gitignore",
            ".ggitignore",
            ".emacs.d/init.el",
            ".skip_on_os",
            ".map_to_os_name",
            ".skip_on_another_os",
        ],
    )
    profile = Profile.load(str(directory))
    links = profile.link_manager.generate_links()
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
            "src": os.path.join(directory, ".emacs.d", "init.el"),
            "dst": os.path.join(os.getenv("HOME"), ".emacs.d", "init.el"),
        },
        {
            "src": os.path.join(directory, ".map_to_os_name"),
            "dst": os.path.join(
                os.getenv("HOME"), ".{name}".format(name=platform.system())
            ),
        },
        {
            "src": os.path.join(directory, ".skip_on_another_os"),
            "dst": os.path.join(os.getenv("HOME"), ".skip_on_another_os"),
        },
    ]

    sorted_links = sorted(links, key=itemgetter("src"))
    sorted_expected_links = sorted(expected_links, key=itemgetter("src"))
    assert len(sorted_links) == len(sorted_expected_links)

    for (link, expected) in zip(sorted_links, sorted_expected_links):
        assert link["src"] == expected["src"]
        assert link["dst"] == expected["dst"]


def test_list_files(dotfile_dir):
    """Test that a profile properly lists it's directory."""
    dotfiles, directory = dotfile_dir()
    link_manager = LinkManager(str(directory))
    assert sorted(link_manager.find_files()) == sorted(dotfiles)


def test_linking(dotdfm):
    """Test that the profile properly creates the links."""
    _, directory = dotdfm()
    with TemporaryDirectory() as target:
        profile = Profile.load(str(directory), extras={"target_dir": target})
        links = profile.link_manager.generate_links()
        profile.link()
        # Use a set comprehension since the target would not contain duplicates
        # since they would be overwritten with the last occuring link.
        dest = sorted(list({os.path.basename(x["dst"]) for x in links}))
        created = []
        for root, dirs, files in os.walk(target):
            dirs[:] = [d for d in dirs if d != ".git"]
            created += [f for f in files if os.path.islink(os.path.join(root, f))]

        created = sorted(created)
        assert dest == created


def test_xdg_config(dotfile_dir):
    _, d_dir = dotfile_dir([".config/nvim/init.vim"])
    profile = Profile.load(str(d_dir))
    links = list(profile.link_manager.generate_links())
    assert links == [
        {
            "src": os.path.join(d_dir, ".config", "nvim", "init.vim"),
            "dst": os.path.join(xdg_dir(), "nvim", "init.vim"),
        },
    ]
