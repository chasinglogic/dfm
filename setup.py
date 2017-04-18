"""dfm, a dotfile manager for pair programmers."""
import os.path as path
import re
from setuptools import find_packages, setup

dependencies = ['click']

dfmfile = path.join(path.dirname(__file__), 'dfm', 'cli.py')
print(dfmfile)

# Thanks to SQLAlchemy:
# https://github.com/zzzeek/sqlalchemy/blob/master/setup.py#L104
with open(dfmfile) as stream:
    __version__ = re.compile(
        r".*__version__ = '(.*?)'", re.S
    ).match(stream.read()).group(1)


setup(
    name='dfm',
    version=__version__,
    url='https://github.com/chasinglogic/dfm',
    download_url='https://github.com/chasinglogic/dfm/tarball/'+__version__,
    license='GPLv3',
    author='Mathew Robinson',
    author_email='chasinglogic@gmail.com',
    description='A dotfile manager for lazy people and pair programmers.',
    packages=find_packages(exclude=['tests']),
    include_package_data=True,
    install_requires=dependencies,
    entry_points={
        'console_scripts': [
            'dfm = dfm.cli:dfm',
        ],
    },
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Environment :: Console',
        'Intended Audience :: Developers',
        'Intended Audience :: System Administrators',
        'License :: OSI Approved :: GNU General Public License v3 (GPLv3)',
        'Operating System :: POSIX',
        'Operating System :: MacOS',
        'Operating System :: Unix',
        'Programming Language :: Python :: 3',
        'Topic :: Software Development :: Libraries :: Python Modules',
        'Topic :: System :: Archiving :: Backup',
        'Topic :: System :: Shells',
        'Topic :: Text Editors',
        'Topic :: Terminals',
        'Topic :: System :: Recovery Tools',
        'Topic :: Utilities',
    ]
)
