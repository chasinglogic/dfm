"""Usage: dfm list

Lists currently downloaded and available profiles.
"""

import os

from dfm.dotfile import dfm_dir


def run(_args):
    for profile in os.listdir(os.path.join(dfm_dir(), 'profiles')):
        if not profile.startswith('.'):
            print(profile)
