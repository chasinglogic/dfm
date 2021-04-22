"""A task management tool that integrates with 3rd party services."""

from setuptools import find_packages, setup


with open("README.md") as f:
    LONG_DESCRIPTION = f.read()

with open("requirements.txt") as f:
    install_requires = f.read().split("\n")

setup(
    name="dfm",
    version="8.4.2",
    url="https://github.com/chasinglogic/dfm",
    license="GPL-3.0",
    author="Mathew Robinson",
    author_email="chasinglogic@gmail.com",
    description="A dotfile manager for pair-programmers and lazy people.",
    long_description=LONG_DESCRIPTION,
    long_description_content_type="text/markdown",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    include_package_data=True,
    zip_safe=False,
    platforms="any",
    install_requires=install_requires,
    entry_points={"console_scripts": ["dfm = dfm.cli:main"]},
    classifiers=[
        # As from https://pypi.org/classifiers/
        "Development Status :: 4 - Beta",
        # 'Development Status :: 5 - Production/Stable',
        "Environment :: Console",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: GNU Affero General Public License v3",
        "Operating System :: POSIX",
        "Operating System :: MacOS",
        "Operating System :: Microsoft :: Windows",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
)
