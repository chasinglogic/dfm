"""Usage: dfm run_hook <hook>

Runs <hook> without the need to invoke the side effects of the given action.
"""

import os
import subprocess
import sys
import yaml

from dfm.cli.utils import current_profile


def run(args):
    dotdfm = os.path.join(current_profile(), '.dfm.yml')
    with open(dotdfm) as dfm_cfg:
        cfg = yaml.load(dfm_cfg)

    hooks = cfg.get('hooks', {})
    commands = hooks.get(args['<hook>'], [])

    for command in commands:
        print('Running script:', command)
        subprocess.run(['/bin/sh', '-c', command],
                       cwd=current_profile(),
                       stdin=sys.stdin,
                       stdout=sys.stdout,
                       stderr=sys.stderr)
