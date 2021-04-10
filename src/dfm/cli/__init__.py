"""
Usage:
    dfm [options] <command> [<args>...]
    dfm help
    dfm sync
    dfm link <profile>

Dotfile management written for pair programmers. Examples on getting
started with dfm are avialable at https://github.com/chasinglogic/dfm

Options:
    -v, --verbose  If provided print more logging info
    --debug        If provided print debug level logging info
    -h, --help     Print this help information

Commands:
    help           Print usage information about dfm commands
    sync (s)       Sync your dotfiles
    add (a)        Add the file to the current dotfile profile
    clean (x)      Clean dead symlinks
    clone (c)      Use git clone to download an existing profile
    git (g)        Run the given git command on the current profile
    init (i)       Create a new profile
    link (l)       Create links for a profile
    list (ls)      List available profiles
    remove (rm)    Remove a profile
    run-hook (rh)  Run dfm hooks without using normal commands
    where (w)      Prints the location of the current dotfile profile

See 'dfm help <command>' for more information on a specific command.
"""

import logging
import sys
from importlib import import_module

from docopt import docopt

from dfm import __version__

ALIASES = {
    "s": "sync",
    "a": "add",
    "x": "clean",
    "c": "clone",
    "g": "git",
    "i": "init",
    "l": "link",
    "ls": "list",
    "rm": "remove",
    "run-hook": "run_hook",
    "rh": "run_hook",
    "w": "where",
}


def main():
    """CLI entrypoint, handles subcommand parsing"""
    args = docopt(
        __doc__,
        version="dfm version {}".format(__version__),
        options_first=True,
    )
    if not args["<command>"]:
        print(__doc__)
        sys.exit(1)

    if args["--debug"]:
        logging.basicConfig(
            stream=sys.stdout,
            level=logging.DEBUG,
        )
    elif args["--verbose"]:
        logging.basicConfig(
            stream=sys.stdout,
            level=logging.INFO,
        )

    logger = logging.getLogger(__name__)

    command = args["<command>"]
    try:
        if command == "help":
            if args["<args>"]:
                help_cmd = ALIASES.get(args["<args>"][0], args["<args>"][0])
                command_mod = import_module("dfm.cli.{}_cmd".format(help_cmd))
                print(command_mod.__doc__)
            else:
                print(__doc__)
            sys.exit(0)

        command = ALIASES.get(command, command)
        command_mod = import_module("dfm.cli.{}_cmd".format(command))
        argv = [command] + args["<args>"]
        command_mod.run(docopt(command_mod.__doc__, argv=argv))
        sys.exit(0)
    except ImportError as exc:
        print("{} is not a known dfm command.".format(command))
        logger.debug("Exception: %s", exc)
        sys.exit(1)
    except KeyboardInterrupt:
        sys.exit(1)


if __name__ == "__main__":
    main()
