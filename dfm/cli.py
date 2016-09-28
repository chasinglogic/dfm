import click
import json
import os
from dfm.lib import get_default_config_dir

CONFIG_DIR = get_default_config_dir()
CONFIG = json.loads(os.path.join(CONFIG_DIR, "config.json"))

@click.group()
@click.option("--config", "-c", 
              type=click.Path(resolve_path=True))
def dfm(config):
    """A dotfile manager for pair programmers."""
    if config != None:
        CONFIG_DIR = config
    pass

@dfm.command()
@click.argument("profile")
def pull(profile):
    """Update a profile."""
    profile_path = get_profile_path(CONFIG_DIR, profile)
    update_profile(profile_path)

@dfm.command()
@click.argument("profile")
def push(profile):
    """Push local changes to the remote."""
    pass
    
@dfm.command()
@click.argument("repo")
def clone(repo):
    """Clone a profile from a git repo."""
    repo_url = get_repo_url(repo)
    profile_path = get_profile_path(CONFIG_DIR, repo)
    clone_profile(repo_url)

@dfm.command()
@click.argument("profile")
def init(profile):
    """Create an empty profile with the given name."""
    pass
    
