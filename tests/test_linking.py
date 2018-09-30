"""Test linking"""

import os
from operator import itemgetter
from tempfile import TemporaryDirectory

from dfm.dotfile import Profile


def setup_module():
    """Set DFM_CONFIG_DIR for this module's tests."""
    os.environ['DFM_CONFIG_DIR'] = TemporaryDirectory().name


def test_translation(dotdfm):
    """Test that the generated links properly translate names."""
    _, directory = dotdfm()
    profile = Profile(directory)
    links = profile.link(dry_run=True)
    expected_links = [{
        'src': os.path.join(directory, 'vimrc'),
        'dst': os.path.join(os.getenv('HOME'), '.vimrc'),
    }, {
        'src': os.path.join(directory, 'bashrc'),
        'dst': os.path.join(os.getenv('HOME'), '.bashrc')
    }, {
        'src': os.path.join(directory, 'emacs'),
        'dst': os.path.join(os.getenv('HOME'), '.emacs')
    }, {
        'src': os.path.join(directory, 'gitignore'),
        'dst': os.path.join(os.getenv('HOME'), '.gitignore')
    }, {
        'src': os.path.join(directory, '.ggitignore'),
        'dst': os.path.join(os.getenv('HOME'), '.gitignore')
    }, {
        'src': os.path.join(directory, 'emacs.d'),
        'dst': os.path.join(os.getenv('HOME'), '.emacs.d'),
        'target_is_directory': True
    }]

    for (link, expected) in zip(
            sorted(links, key=itemgetter('src')),
            sorted(expected_links, key=itemgetter('src'))):
        assert link['src'] == expected['src']
        assert link['dst'] == expected['dst']
        assert link.get('target_is_directory') == expected.get(
            'target_is_directory')


def test_linking(dotdfm):
    """Test that the profile properly creates the links."""
    _, directory = dotdfm()
    with TemporaryDirectory() as target:
        profile = Profile(directory, target_dir=target)
        links = profile.link()
        # Use a set comprehension since the target would not contain duplicates
        # since they would be overwritten with the last occuring link.
        dest = sorted(list({os.path.basename(x['dst']) for x in links}))
        created = sorted(os.listdir(target))
        print(created)
        print(dest)
        assert dest == created
