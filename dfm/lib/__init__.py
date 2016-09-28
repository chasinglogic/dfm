import os
import os.path as path
from subprocess import run, PIPE

def get_default_config_dir():
    home = os.environ.get("HOME", 
                          os.environ.get("APPDATA"),
                          "")
    return path.join(home, ".config", "dfm")

def get_repo_url(repo_name):
    if len(repo_name.split("/")) == 2:
        return "https://github.com/" + repo_name
    return repo_name

def get_profile_path(config_dir, profile_name):
    spl = profile_name.split("/")
    if len(spl) > 1:
        # if we got passed a url of some sort return
        # the second to last element of the split.
        return path.join(config_dir, spl[len(spl) - 3])
    return path.join(config_dir, profile_name)

def update_profile(path):
    run([ "git", "pull", "origin", "master" ], 
        cwd=path, stdout=PIPE)

def clone_profile(repo, profile_path):
    run([ "git", "clone", repo, profile_path ], stdout=PIPE)
