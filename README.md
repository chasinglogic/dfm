# dfm

A dotfile manager for lazy people and pair programmers.

> dfm doesn't require that the dotfiles in your repository start with dots
> though it handles either case equally well.

## Table of Contents

- [Features](#features)
  - [Multiple dotfile profiles](#multiple-dotfile-profiles)
  - [Profile modules](#profile-modules)
  - [Pre and post command hooks](#pre-and-post-command-hooks)
  - [Respects `$XDG_CONFIG_HOME`](#respects-xdg_config_home)
  - [Skips relevant files](#skips-relevant-files)
  - [Configurable mappings](#file-mappings)
- [Installation](#installation)
- [Updating](#updating)
- [Usage](#usage)
- [Git quick start](#git-quick-start)
  - [Existing dotfiles repository](#quick-start-existing-dotfiles-repository)
  - [No existing dotfiles repository](#quick-start-no-existing-dotfiles-repository)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features

dfm supports these features that I was unable to find in other Dotfile
Management solutions.

### Multiple dotfile profiles

dfm's core feature is the idea of profiles. Profiles are simply a
collection of dotfiles that dfm manages and links in the `$HOME`
directory or configuration directories. This means that you can have
multiple profiles and overlap them.

This feature is hard to describe, so I will illustrate it's usefulness
with two use cases:

#### The work profile

I use one laptop for work and personal projects in my dfm profiles I have my
personal profile `chasinglogic` which contains all my dotfiles for Emacs, git,
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

Since dfm when linking only overwrites the files which are in the new
profile, I can run `dfm link work` and still have access to my emacs
configuration but my `gitconfig` has been updated to use my work
email. Similarly when I leave work I just `dfm link chasinglogic` to
switch back.

See [profile modules](#profile-modules) for an even better solution to this
particular use case.

#### Pair programming

The original inspiration for this tool was pair programming with my
friend [lionize](https://github.com/lionize). lionize has a dotfiles
repository so I can clone it using the git backend for dfm with `dfm
clone --name lionize https://github.com/lionize/dotfiles`.

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

### Profile modules

dfm supports profile modules which can be either additional dotfiles profiles as
accepted by the `dfm clone` command or can be any git repository such as
[Spacemacs](https://github.com/syl20bnr/spacemacs). You can get more info about
how to use them and configure them in [Configuration](#configuration)

### Pre and Post command hooks

dfm supports pre and post command hooks. that allows you to specify before and
after command shell scripts to run. For example, I use a profile module to keep
certain ssh keys in an encrypted git repository. Whenever I run the `dfm sync` command
I have hooks which fix the permissions of the keys and ssh-add them to my ssh
agent. You can read about how to write your own hooks in
[Configuration](#configuration)

### Respects `$XDG_CONFIG_HOME`

dfm respects dotfiles which exist in the `$XDG_CONFIG_HOME` directory,
meaning if in your repository you have a folder named `config` or
`.config` it'll translate those into the `$XDG_CONFIG_HOME`
directory automatically. Similarly when using `dfm add` if inside your
`$XDG_CONFIG_HOME` or $HOME/.configuration directories it'll add those to
the repository appropriately.

### Skips relevant files

dfm by default will skip multiple relevant files.

- .git

dfm will skip .git directory your `$HOME` directory isn't turned into
a git repository.

- .gitignore

If you would like to store a global `.gitignore` file you can either omit the
leading dot (so just `gitignore`) or name the global one `.ggitignore` and dfm
will translate the name for you. Otherwise it assumes that `.gitignore` is the
gitignore for the profile's repository and so skips it.

- README

Want to make a README for your dotfiles? Go ahead! As long as the file name
starts with README dfm will ignore it. So `README.txt` `README.md` and
`README.rst` or whatever other permutations you can dream up all work.

- LICENSE

You should put a LICENSE on all code you put on the internet and some dotfiles /
configurations are actual code (See: Emacs). If you put a LICENSE in your
profile dfm will respect you being a good open source citizen and not clutter your
`$HOME` directory.

- .dfm.yml

This is a special dfm file used for hooks today and in the future for other ways
to extend dfm. As such dfm doesn't put it in your `$HOME` directory.

### Custom mappings

The above ignores are implemented as a dfm feature called Mappings. You can
write your own mappings to either skip or translate files to different
locations than dfm would normally place them. You can read how to configure your
own mappings in [Configuration](#configuration)

## Installation

### Install from Release

dfm is available on PyPi and should be installed from there:

```text
$ pip3 install dfm
```

dfm supports Python 3+.

### Install from Source

Clone the repository and run `make install`:

```bash
git clone https://github.com/chasinglogic/dfm
cd dfm
make install
```

> It's possible that for your system you will need to run the make
> command with sudo.

## Usage

```text
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
```

## Quick start

### Quick start (Existing dotfiles repository)

If you already have a dotfiles repository you can start by cloning it using the clone
command.

> SSH URLs will work as well.

```bash
dfm clone https://github.com/chasinglogic/dotfiles
```

If you're using GitHub you can shortcut the domain:

```bash
dfm clone chasinglogic/dotfiles
```

If you want to clone and link the dotfiles in one command:

```bash
dfm clone --link chasinglogic/dotfiles
```

You may have to use `--overwrite` as well if you have existing non-symlinked
versions of your dotfiles

Once you have multiple profiles you can switch between them using `dfm link`

```bash
dfm link some-other-profile
```

See the Usage Notes below for some quick info on what to expect from other dfm
commands.

### Quick Start (No existing dotfiles repository)

If you don't have a dotfiles repository the best place to start is with `dfm init`

```bash
dfm init my-new-profile
```

Then run `dfm link` to set it as the active profile, this is also how you switch
profiles

```bash
dfm link my-new-profile
```

Once that's done you can start adding your dotfiles

```bash
dfm add ~/.bashrc
```

Alternatively you can add multiple files at once

```bash
dfm add ~/.bashrc ~/.vimrc ~/.vim ~/.emacs.d
```

Then create your dotfiles repository on GitHub. Instructions for how to do that can be
found [here](https://help.github.com/articles/create-a-repository/). Once that's done
get the "clone" URL for your new repository and set it as origin for the profile:

**Note:** When creating the remote repository don't choose any options such as
"initialize this repository with a README" otherwise git'll get cranky when you add
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

## Configuration

dfm supports a `.dfm.yml` file in the root of your repository that
changes dfm's behavior when syncing and linking your profile. This
file will be ignored when doing a `dfm link` so won't end up in
your home directory. The `.dfm.yml` can be used to configure these
features:

- [Modules](#modules)
- [Mappings](#mappings)
- [Hooks](#hooks)

### Modules

Modules in dfm are like sub profiles. They're git repositories that are cloned into a
a special directory: `$XDG_CONFIG_HOME/dfm/modules`. They're shared across
profiles so if two dotfile profiles have the same module they'll share that
module.

The syntax for defining a minimum module is as follows:

```yaml
modules:
    - repository: git@github.com:chasinglogic/dotfiles
```

This would clone my dotfiles repository as a module into
`$XDG_CONFIG_HOME/dfm/modules/dotfiles`. If I wanted to use a unique name or
some other folder name so it wouldn't be shared you can specify an additional
option `name`:

```yaml
modules:
    - repository: git@github.com:chasinglogic/dotfiles
      name: chasinglogic-dotfiles
```

Which would instead clone into
`$XDG_CONFIG_HOME/dfm/modules/chasinglogic-dotfiles`. You can define multiple
modules:

```yaml
modules:
    - repository: git@github.com:chasinglogic/dotfiles
      name: chasinglogic-dotfiles
    - repository: git@github.com:lionize/dotfiles
```

Make sure that you specify a name if the resulting clone location as defined by
git would conflict as we see here. Both of these would have been cloned into
dotfiles which would cause the clone to fail for the second module if we didn't
specify name for chasinglogic's dotfiles.

An additional use for modules is that of a git repository you want to clone but not
link. An example use would be for downloading
[Spacemacs](https://github.com/syl20bnr/spacemacs) or any such community
configuration like oh-my-zsh, etc.

```yaml
modules:
    - repository: git@github.com:syl20bnr/spacemacs
      link: none
      pull_only: true
      location: ~/.emacs.d
```

Here we specify a few extra keys. There purpose should be self explanatory but
if you're curious [below](#available-keys) is a detailed explanation of all keys
that each module configuration supports.

#### Available keys

- [repository](#repository)
- [name](#name)
- [location](#location)
- [link](#link)
- [pull\_only](#pull\_only)
- [mappings](#mappings)

##### repository

Required, this is the git repository to clone for the module.

##### name

This changes the cloned name. This only has an effect if location isn't
provided. Normally a git repository would be cloned into
`$XDG_CONFIG_HOME/dfm/modules` and the resulting folder would be named whatever
git decides it should be based on the git URL. If this is provided it'll be
cloned into the modules directory with the specified name. This is useful if
multiple profiles use the same module.

##### location

If provided module will be cloned into the specified location. You can use the
`~` bash expansion here to represent `$HOME`. No other expansions are available.
This option is useful for cloning community configurations like oh-my-zsh or
spacemacs.

##### link

Determines when to link the module. Link in this context means that it'll be
treated like a normal dotfile profile, so all files will go through the same
translation rules as a regular profile and be linked accordingly. Available
values are `post`, `pre`, and `none`. `post` is the default and means that the
module will be linked after the parent profile. "pre" means this will be linked
before the parent profile, use this if for instance you want to use most files
from this profile and override a few files with those from the parent file since
dfm will overwrite the links with the last one found. "none" means the module is
not a dotfiles profile and shouldn't be linked at all, an example being
community configuration repositories like oh-my-zsh or spacemacs.

##### pull\_only

If set to `true` won't attempt to push any changes. It's important to
know that dfm always tries to push to origin master, so if you don't
have write access to the repository or don't want it to automatically
push to master then you should set this to true. This is useful for
community configuration repositories.

##### mappings

A list of file mappings as described below in [Mappings](#mappings). Modules do
not inherit parent mappings, they do however inherit the default mappings as
described in [Skips Relevant Files](#skips-relevant-files)

### Mappings

Mappings are a way of defining custom file locations. To understand
mappings one must understand dfm's default link behavior:

#### Default behavior

For an example let's say you have a file named `my_config.txt` in your
dotfile repository. dfm will try and translate that to a new location
of `$HOME/.my_config.txt`. It'll then create a symlink at that location
pointing to `my_config.txt` in your dotfile repository.

#### Using mappings

With mappings you can replace this behavior and make it so dfm will
link `my_config` wherever you wish. This is useful if you need to
store config files that are actually global. Such as configuration
files that would go into `/etc/` or if you want to sync some files in
your repo but not link them.

Here is a simple example:

```yaml
mappings:
  - match: my_global_etc_files
    target_dir: /etc/
  - match: something_want_to_skip_but_sync
    skip: true
```

Here dfm uses the match as a regular expression to match the file
paths in your dotfile repository. When it finds a path which matches
the regular expression it adds an alternative linking behavior. For
anything where `skip` is true it simply skips linking. For anything
with `target_dir` that value will override `$HOME` when linking.

#### Available configuration

Mappings support the following configuration options:

- [match](#match)
- [skip](#skip)
- [target\_dir](#target\_dir)

##### match

Match is a regular expression used to match the file path of any files
in your dotfile repository. This is used to determine if the custom
linking behavior for a file should be used.

These are python style regular expressions and are matched using the
[`re.findall`](https://docs.python.org/3/library/re.html#re.findall)
method so are by default fuzzy matching.

##### skip

If provided the file/s will not be linked.

##### target_dir

Where to link the file to. The `~` expansion for `$HOME` is supported
here but no other expansions are available. It is worth noting that if
you're using `~` in your target_dir then you should probably just
create the directory structure in your git repo.

### Hooks

Hooks in dfm are used for those few extra tasks that you need to do whenever
your dotfiles are synced or linked.

An example from my personal dotfiles is running an Ansible playbook
whenever I sync my dotfiles. To accomplish this I wrote an
`after_sync` hook as follows:

```yaml
hooks:
  after_sync:
    - ansible-playbook ansible/dev-mac.yml
```

Now whenever I sync my dotfiles Ansible will run my `dev-mac` playbook to make
sure that my packages etc are also in sync!

The hooks option is just a YAML map which supports the following keys:
`after_link`, `before_link`, `after_sync`, and `before_sync`. The
values of any of those keys is a YAML list of strings which will be
executed in a shell via `/bin/sh -c '$YOUR COMMAND'`. An example would
be:

```yaml
hooks:
  after_link:
    - ls -l
    - whoami
    - echo "All done!"
```

All commands are ran with a working directory of your dotfile
repository and the current process environment is passed down to the
process so you can use `$HOME` etc environment variables in your
commands.

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
    Copyright (C) 2018 Mathew Robinson

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
