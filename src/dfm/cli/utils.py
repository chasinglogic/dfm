"""CLI utility functions used throughout dfm"""

import json
import os
import sys

from dfm.dotfile import Profile, dfm_dir


def state_file_p():
    """Return the path to the dfm state file."""
    return os.path.join(dfm_dir(), 'state.json')


def load_profile(name):
    """
    Load a profile by name.

    Joins the dfm state directory with 'profiles' and name to
    determine where the profile is.
    """
    path = os.path.join(dfm_dir(), 'profiles', name)
    return Profile(path)


def switch_profile(name):
    """
    Switch profile will update the state file to the profile with name.

    Returns the profile object as returned by load_profile.
    """
    path = state_file_p()
    with open(path, 'w+') as state_file:
        content = state_file.read()
        if content:
            state = json.loads(content)
        else:
            state = {}

        state['current_profile'] = name
        json.dump(state, state_file)
        return load_profile(name)


def current_profile():
    """
    Load the current profile as indicated in the state file.

    Exits with a helpfule error message if state file is not found or
    current_profile is not set.
    """
    try:
        with open(state_file_p()) as state_file:
            state = json.load(state_file)
    except FileNotFoundError:
        state = {}

    profile_name = state.get('current_profile')
    if not profile_name:
        print('no profile active, run dfm link to make one active')
        sys.exit(1)

    return os.path.join(dfm_dir(), 'profiles', profile_name)


def inject_profile(wrapped):
    """Inject the current profile as a keyword argument 'profile'."""

    def wrapper(*args, **kwargs):
        kwargs['profile'] = Profile(current_profile())
        return wrapped(*args, **kwargs)

    return wrapper
