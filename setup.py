"""
A dotfile manager for pair programmers.
"""
from setuptools import find_packages, setup

dependencies = ['click']

setup(
    name='dfm',
    version='0.3',
    url='https://github.com/chasinglogic/dfm',
    download_url='https://github.com/chasinglogic/dfm/tarball/0.3',
    license='GPLv3',
    author='Mathew Robinson',
    author_email='mathew.robinson3114@gmail.com',
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
