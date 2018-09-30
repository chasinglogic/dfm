"""CLI utility functions used throughout dfm"""

import json
import sys
from os.path import join

from dfm.dotfile import Profile, dfm_dir


def state_file_p():
    return join(dfm_dir(), 'state.json')


def load_profile(name):
    path = join(dfm_dir(), 'profiles', name)
    return Profile(path)


def switch_profile(name):
    with open(state_file_p(), 'w+') as state_file:
        state = json.load(state_file)
        state['current_profile'] = name
        return load_profile(name)


def current_profile():
    with open(state_file_p()) as state_file:
        state = json.load(state_file)
        profile = state.get('current_profile')
        if profile is None:
            print('no profile active, run dfm link to make one active')
            sys.exit(1)
        return join(dfm_dir(), 'profiles', profile)


def inject_profile(wrapped):
    """Inject the current profile as a keyword argument 'profile'."""

    def wrapper(*args, **kwargs):
        kwargs['profile'] = Profile(current_profile())
        return wrapped(*args, **kwargs)

    return wrapper
