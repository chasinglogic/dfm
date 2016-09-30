#
#     dfm, a dotfile manager for lazy people and pair programmers
#
#     Copyright (C) 2016 Mathew Robinson <mathew.robinson3114@gmail.com>
#
#     This program is free software: you can redistribute it and/or modify
#     it under the terms of the GNU General Public License as published by
#     the Free Software Foundation, either version 3 of the License, or
#     (at your option) any later version.
#
#     This program is distributed in the hope that it will be useful,
#     but WITHOUT ANY WARRANTY; without even the implied warranty of
#     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#     GNU General Public License for more details.
#
#     You should have received a copy of the GNU General Public License
#     along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
import os
from pathlib import Path
import shutil
import click
from dfm.config import CONFIG
from subprocess import run, PIPE

def get_repo_url(repo_name):
    if len(repo_name.split("/")) == 2:
        return "https://github.com/" + repo_name
    return repo_name

def get_profile_path(config_dir, profile_name):
    spl = profile_name.split("/")
    if len(spl) > 1:
        # if we got passed a url of some sort return
        # the second to last element of the split.
        return os.path.join(config_dir, "profiles", spl[len(spl) - 2])
    return os.path.join(config_dir, "profiles", profile_name)

def pull_profile(path, branch):
    run([ "git", "pull", "origin", branch ],
        cwd=path, stdout=PIPE, stdin=PIPE)

def clone_profile(repo, profile_path):
    run([ "git", "clone", repo, profile_path ],
        stdout=PIPE, stdin=PIPE)

def push_profile(path, branch):
    run([ "git", "push", "origin", branch ],
        cwd=path, stdout=PIPE, stdin=PIPE)

def create_and_init_profile(profile_path):
    click.echo("Creating profile %s" % profile_path)
    os.mkdir(profile_path)
    run([ "git", "init" ],
        cwd=profile_path,
        stdout=PIPE)

def gen_dot_file(flname):
    df = flname if flname.startswith(".") else "." + flname
    return os.path.abspath(os.path.join(os.environ.get("HOME", ""), df))


def link_file(fl, df, force=False):
    # Check if a non sym linked version exists.
    if ((os.path.exists(df) and
         not os.path.islink(df)) and not force):
        click.echo("Error linking: %s" % df)
        click.echo("Dotfile exists and isn't symlink. Refusing to overwrite. Use --force to overwrite")
        return

    if os.path.isfile(df) or os.path.islink(df):
        os.remove(df)
    if os.path.isdir(df):
        shutil.rmtree(df)

    os.symlink(fl, df)
    if CONFIG["verbose"]:
        click.echo("Linked file %s -> %s" % (fl, df))

def link_profile(profile_path, force=False):
    dot_files = os.scandir(os.path.abspath(profile_path))
    click.echo("Linking profile %s" % profile_path)
    for d in dot_files:
        # Skip the git directory
        if d.name == ".git":
            continue

        if d.name == "config" or d.name == ".config":
            # Get the .config directory files
            dfgf = os.scandir(d.path)
            for f in dfgf:
                # .config files have a different path
                config_path = os.environ.get("XDG_CONFIG_HOME", "")
                if xdg == "":
                    config_path = os.path.join(os.environ.get("HOME", ""), ".config")

                dfp = os.path.join(config_path, f.name)
                link_file(f.path, dfp, force=force)
            continue
        link_file(d.path, gen_dot_file(d.name), force=force)

def add_file(path, profile):
    xdg = os.environ.get("XDG_CONFIG_HOME", "")
    old_file = Path(path)

    fn = old_file.name
    # If starts with a dot remove it.
    if fn.startswith("."):
        fn = fn.replace(".", "", 1)

    new_path = os.path.join(profile, fn)
    if Path(xdg) in { old_file }:
        new_path = os.path.join(profile, "config", fn)

    os.rename(bytes(old_file), new_path)
    link_file(new_path, bytes(old_file), force=True)
    run([ "git", "add", fn ], cwd=profile, stdout=PIPE)

def checkout_profile(profile, branch):
    run([ "git", "checkout", branch ],
        cwd=profile, stdout=PIPE)

def commit_profile(profile, message):
    run([ "git", "commit", "-am", message ],
        cwd=profile, stdout=PIPE)

def set_remote_profile(profile, remote):
    run([ "git", "remote", "set-url", "origin", remote ],
        cwd=profile, stdout=PIPE)
