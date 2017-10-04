# dfm
A dotfile manager for lazy people and pair programmers.

**NOTE:** DFM does not require that the dotfiles in your repo start with dots
though it handles either case equally well.

## Table of Contents

- [Features](#features)
  - [Multiple Dotfile Profiles](#multiple-dotfile-profiles)
  - [Pre and Post Command Hooks](#pre-and-post-command-hooks)
  - [Respects $XDG\_CONFIG\_HOME](#respects-xdg_config_home)
  - [Skips Relevant Files](#skips-relevant-files)
- [Installation](#installation)
- [Updating](#updating)
- [Usage](#usage)
- [Git Quick Start](#git-quick-start)
  - [Existing Dotfiles Repo](#quick-start-existing-dotfiles-repo)
  - [No Existing Dotfiles Repo](#quick-start-no-existing-dotfiles-repo)
- [Contributing](#contributing)
- [License](#license)

## Features

DFM Supports these features that I was unable to find in other Dotfile
Management solutions.

### Multiple Dotfile Profiles

DFM's core feature is the idea of "profiles". Profiles are simply a collection
of dotfiles that DFM manages and links in the `$HOME` directory or config
directories. This means that you can have multiple profiles and overlap them.
This feature is hard to write directly about so I will illustrate it's
usefulness with two Use Cases:

#### The Work Profile

I use one laptop for work and personal projects in my DFM profiles I have my
personal profile "chasinglogic" which contains all my dotfiles for Emacs, git,
etc. and a "work" profile which only has a `.gitconfig` that has my work email
in it. So my profile directory looks like this:

```text
profiles/
├── chasinglogic
│   ├── agignore
│   ├── bash
│   ├── bashrc
│   ├── gitconfig
│   ├── gnupg
│   ├── password-store
│   ├── pypirc
│   ├── spacemacs.d
│   └── tmux.conf
└── work
    └── gitconfig
```

Sinc dfm when linking only does the deltas I can run `dfm link work` and still
have access to my emacs config but my gitconfig has been updated to use my work
email. Simliarly when I leave work I just `dfm link chasinglogic` to switch back.

#### Pair Programming

The original inspiration for this tool was pair programming with my friend
[lionize](https://github.com/lionize). lionize has a dotfiles repo so I can
clone it using the git backend for DFM with `dfm clone lionize/dotfiles`. Note
that if a partial URL like this one is given DFM will assume you want to clone
via https from Github but full git-cloneable URLs (including SSH) can be passed
to this command.

Now our profile directory looks like:

```text
profiles/
├── chasinglogic
│   ├── .dfm.yml
│   ├── .git
│   ├── .gitignore
│   ├── agignore
│   ├── bash
│   ├── bashrc
│   ├── gitconfig
│   ├── gnupg
│   ├── password-store
│   ├── pypirc
│   ├── spacemacs.d
│   └── tmux.conf
├── lionize
│   ├── .agignore
│   ├── .git
│   ├── .gitconfig
│   ├── .gitignore_global
│   ├── .gitmessage
│   ├── .scripts
│   ├── .tmux.conf
│   ├── .vim
│   ├── .vimrc -> ./.vim/init.vim
│   └── .zshrc
└── work
    ├── .git
    └── gitconfig
```

Now when I'm driving I simply `dfm link chasinglogic` and when passing back to
lionize he runs `dfm link lionize` and we don't have to mess with multiple
machines vice versa.

### Pre and Post command hooks

DFM supports a `.dfm.yml` file in the root of your repo that has a "hooks"
key that allows you to specify before and after command bash scripts to run.
For example, I use Spacemacs so whenever I run `dfm clone chasinglogic/dotfiles`
I want it to install spacemacs. Additionally, I have it set up to update
Spacemacs whenever I run `dfm sync`. Here is my `.dfm.yml`:

```yaml
hooks:
  after_sync:
    - |
      cd ~/.emacs.d
      git pull
      echo "Updated Spacemacs!"
  after_clone:
    - |
      if [ ! -d ~/.emacs.d ]; then
          git clone https://github.com/syl20bnr/spacemacs ~/.emacs.d
      fi
```

You can wrap any dfm subcommand, the syntax is simply `before_{command_name}` or
`after_{command_name}` and then a yaml list of bash commands to run.

### Respects ``XDG\_CONFIG\_HOME`

dfm respects dotfiles which exist in the $XDG\_CONFIG\_HOME directory, meaning
if in your repo you have a folder named config or .config it will translate
those into the $XDG\_CONFIG\_HOME directory appropriately. Similarly when
using `dfm add` if inside your $XDG\_CONFIG\_HOME or $HOME/.config directories
it will add those to the repo appropriately.

### Skips relevant files

Of course DFM skips your .git directory but additionally it will skip these
files:

- .gitignore

If you would like to store a global `.gitignore` file you can either omit the
leading dot (so just `gitignore`) or name the global one `.ggitignore` and DFM
will translate the name for you. Otherwise it assumes that `.gitignore` is the
gitignore for the profile's repo and so skips it.

- README

Want to make a README for your dotfiles? Go ahead! As long as the file name
starts with README DFM will ignore it. So `README.txt` `README.md` and
`README.rst` or whatever other permutations you can dream up all work.

- LICENSE

You should put a LICENSE on all code you put on the internet and some dotfiles /
configurations are actual code (See: Emacs). If you put a LICENSE in your
profile DFM will respect you being a good internet citizen and not clutter your
`$HOME` directory.

- .dfm.yml

This is a special DFM file used for hooks today and in the future for other ways
to extend DFM. As such DFM doesn't put it in your `$HOME` directory.

## Installation

### Install from Release

1. Navigate to [the Releases Page](https://github.com/chasinglogic/dfm/releases)
2. Find the tar ball for your platform / architecture. For example, on 64 bit
   Mac OSX, the archive is named `dfm_{version}_darwin_amd64.tar.gz`
3. Extract the tar ball
4. Put the dfm binary in your `$PATH`

### Install from Source

Simply run go get:

```bash
$ go get github.com/chasinglogic/dfm
```

If your `$GOPATH/bin` is in your `$PATH` then you now have dfm installed.

## Updating

`dfm` can update itself to bring in the latest bug fixes and features. Simply run:

```bash
dfm update
```

To update.

## Usage

```text
Dotfile management written for pair programmers. Examples on getting
started with dfm are avialable at https://github.com/chasinglogic/dfm

Usage:
  dfm [command]

Available Commands:
  add         Add a file to the current profile.
  clean       clean dead symlinks
  clone       git clone an existing profile from `URL`
  git         run the given git command on the current profile
  help        Help about any command
  init        create a new profile with `NAME`
  link        link the profile with `NAME`
  list        list available profiles
  remove      remove the profile with `NAME`
  sync        sync the current profile with the configured backend
  update      downlaod and install dfm updates
  version     print version information for dfm
  where       prints the current profile directory path

Flags:
  -d, --dry-run   don't make changes just print what would happen
  -h, --help      help for dfm
  -v, --verbose   verbose output

Use "dfm [command] --help" for more information about a command.
```

## Git Quick Start

DFM supports pluggable backends, so if you don't like Git you can choose another
option, but git is the default backend so here is a Quick Start Guide to get you going!

### Quick Start (Existing dotfiles repo)

If you already have a dotfiles repo you can start by cloning it using the clone
command.

**Note:** ssh urls will work as well.

```bash
dfm clone https://github.com/chasinglogic/dfiles
```

If you're using github you can shortcut the domain:

```bash
dfm clone chasinglogic/dfiles
```

If you want to clone and link the dotfiles in one command:

```bash
dfm clone --link chasinglogic/dfiles
```

You may have to use `--force` as well if you have existing non-symlinked
versions of your dotfiles

Once you have multiple profiles you can switch between them using `dfm link`

```bash
dfm link some-other-profile
```

See the Usage Notes below for some quick info on what to expect from other dfm
commands.

### Quick Start (No existing dotfiles repo)

If you do not have a dotfiles repo the best place to start is with `dfm init`

```bash
dfm init my-new-profile
```

Then run `dfm link` to set it as the active profile, this is also how you switch
profiles

```bash
dfm link my-new-profile
```

Once that is done you can start adding your dotfiles

```bash
dfm add ~/.bashrc
```

Alternatively you can add multiple files at once

```bash
dfm add ~/.bashrc ~/.vimrc ~/.vim ~/.emacs.d
```

Then create your dotfiles repo on Github. Instructions for how to do that can be
found [here](https://help.github.com/articles/create-a-repo/). Once that's done
get the "clone" URL for your new repo and set it as origin for the profile:

**Note:** When creating the remote repo do not choose any options such as
"initialize this repo with a README" otherwise git will get cranky when you add
the remote because of a recent git update and how it handles [unrelated
histories](http://stackoverflow.com/questions/37937984/git-refusing-to-merge-unrelated-histories)
if you do don't worry the linked post explains how to get past it.

```bash
dfm git remote add origin <your clone URL>
```

Then simply run `dfm sync` to sync your dotfiles to the remote

```bash
dfm sync
```

Now you're done!

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. :fire: Submit a pull request :D :fire:

All pull requests should go to the develop branch not master. Thanks!

## License

This code is distributed under the GNU General Public License

```
    Copyright (C) 2017 Mathew Robinson

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
