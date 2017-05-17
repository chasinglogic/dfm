"""Contains functions for loading and running hooks."""

import os
import shlex

from yaml import load
from os.path import exists
from os.path import join
from functools import wraps
from subprocess import run

from dfm.config import CONFIG


def auto_hooks(fn):
    """Run hooks based on the name of fn."""
    @wraps(fn)
    def wrapper(*args, **kwargs):
        path = CONFIG.get('profile', None)
        if path is None:
            return fn(*args, **kwargs)

        hooks = load_hooks(path)
        if hooks is None:
            return fn(*args, **kwargs)

        run_hooks(hooks, 'before_'+fn.__name__, path=path)
        res = fn(*args, **kwargs)
        run_hooks(hooks, 'after_'+fn.__name__, path=path)
        return res
    return wrapper


def run_hooks(hooks, hook_name, path=os.getcwd()):
    """Run hooks indicated by hook_name."""
    hks = hooks.get(hook_name, [])
    for h in hks:
        run(shlex.split(h), cwd=path)


def load_hooks(path):
    """Load all hooks at path."""
    if exists(join(path, '.dfm.yml')):
        yml = join(path, '.dfm.yml')
    elif exists(join(path, '.dfm.yaml')):
        yml = join(path, '.dfm.yaml')
    else:
        return None

    with open(yml) as f:
        loaded = load(f)

    return loaded['hooks']
