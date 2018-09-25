# dfm
A dotfile manager for lazy people and pair programmers.

**NOTE:** dfm does not require that the dotfiles in your repo start with dots
though it handles either case equally well.

## Table of Contents

- [Features](#features)
  - [Multiple Dotfile Profiles](#multiple-dotfile-profiles)
  - [Profile Modules](#profile-modules)
  - [Pre and Post Command Hooks](#pre-and-post-command-hooks)
  - [Respects $XDG\_CONFIG\_HOME](#respects-xdg_config_home)
  - [Skips Relevant Files](#skips-relevant-files)
  - [Configurable File Mappings](#file-mappings)
- [Installation](#installation)
- [Updating](#updating)
- [Usage](#usage)
- [Git Quick Start](#git-quick-start)
  - [Existing Dotfiles Repo](#quick-start-existing-dotfiles-repo)
  - [No Existing Dotfiles Repo](#quick-start-no-existing-dotfiles-repo)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features

dfm supports these features that I was unable to find in other Dotfile
Management solutions.

### Multiple Dotfile Profiles

dfm's core feature is the idea of "profiles". Profiles are simply a collection
of dotfiles that dfm manages and links in the `$HOME` directory or config
directories. This means that you can have multiple profiles and overlap them.
This feature is hard to write directly about so I will illustrate it's
usefulness with two Use Cases:

#### The Work Profile

I use one laptop for work and personal projects in my dfm profiles I have my
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
email. Simliarly when I leave work I just `dfm link chasinglogic` to switch
back.

See [profile modules](#profile-modules) for an even better solution to this
particular use case.

#### Pair Programming

The original inspiration for this tool was pair programming with my friend
[lionize](https://github.com/lionize). lionize has a dotfiles repo so I can
clone it using the git backend for dfm with `dfm clone lionize/dotfiles`. Note
that if a partial URL like this one is given dfm will assume you want to clone
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

### Profile Modules

dfm supports profile modules which can be either additional dotfiles profiles as
accepted by the `dfm clone` command or can be any git repo such as
[Spacemacs](https://github.com/syl20bnr/spacemacs). You can get more info about
how to use them and configure them in [Configuration](#configuration)

### Pre and Post command hooks

dfm supports pre and post command hooks. that allows you to specify before and
after command shell scripts to run. For example, I use a profile module to keep
certain ssh keys in an encrypted git repo. Whenever I run the `dfm sync` command
I have hooks which fix the permissions of the keys and ssh-add them to my ssh
agent. You can read about how to write your own hooks in
[Configuration](#configuration)

### Respects ``XDG\_CONFIG\_HOME`

dfm respects dotfiles which exist in the $XDG\_CONFIG\_HOME directory, meaning
if in your repo you have a folder named config or .config it will translate
those into the $XDG\_CONFIG\_HOME directory appropriately. Similarly when using
`dfm add` if inside your $XDG\_CONFIG\_HOME or $HOME/.config directories it will
add those to the repo appropriately.

### Skips relevant files

Of course dfm skips your .git directory but additionally it will skip these
files:

- .gitignore

If you would like to store a global `.gitignore` file you can either omit the
leading dot (so just `gitignore`) or name the global one `.ggitignore` and dfm
will translate the name for you. Otherwise it assumes that `.gitignore` is the
gitignore for the profile's repo and so skips it.

- README

Want to make a README for your dotfiles? Go ahead! As long as the file name
starts with README dfm will ignore it. So `README.txt` `README.md` and
`README.rst` or whatever other permutations you can dream up all work.

- LICENSE

You should put a LICENSE on all code you put on the internet and some dotfiles /
configurations are actual code (See: Emacs). If you put a LICENSE in your
profile dfm will respect you being a good internet citizen and not clutter your
`$HOME` directory.

- .dfm.yml

This is a special dfm file used for hooks today and in the future for other ways
to extend dfm. As such dfm doesn't put it in your `$HOME` directory.

### File Mappings

The above ignores are implemented as a DFM feature called File Mappings. You can
write your own file mappings to either skip or translate files to different
locations than dfm would normally place them. You can read how to configure your
own mappings in [Configuration](#configuration)

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

`dfm` can update itself to bring in the latest bug fixes and features. Simply
run:

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

dfm supports pluggable backends, so if you don't like Git you can choose another
option, but git is the default backend so here is a Quick Start Guide to get you
going!

### Quick Start (Existing dotfiles repo)

If you already have a dotfiles repo you can start by cloning it using the clone
command.

**Note:** ssh urls will work as well.

```bash
dfm clone https://github.com/chasinglogic/dotfiles
```

If you're using github you can shortcut the domain:

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

## Configuration

dfm supports a `.dfm.yml` file in the root of your repo that provides various
configuration options. This file will be ignored when doing a `dfm link` so will
not end up in your home directory. The `.dfm.yml` can be used to configure these
features:

- [Modules](#modules)
- [Mappings](#mappings)
- [Hooks](#hooks)

### Modules

Modules in dfm are like sub profiles. They are git repos that are cloned into a
a special directory: `$XDG_CONFIG_HOME/dfm/modules`. They are shared across
profiles so if two dotfile profiles have the same module they will share that
module.

The syntax for defining a minimum module is as follows:

```yaml
modules:
    - repo: git@github.com:chasinglogic/dotfiles
```

This would clone my dotfiles repo as a module into
`$XDG_CONFIG_HOME/dfm/modules/dotfiles`. If I wanted to use a unique name or
some other folder name so it wouldn't be shared you can specify an additional
key `name`:

```yaml
modules:
    - repo: git@github.com:chasinglogic/dotfiles
      name: chasinglogic-dotfiles
```

Which would instead clone into 
`$XDG_CONFIG_HOME/dfm/modules/chasinglogic-dotfiles`. You can define multiple
modules:

```yaml
modules:
    - repo: git@github.com:chasinglogic/dotfiles
      name: chasinglogic-dotfiles
    - repo: git@github.com:lionize/dotfiles
```

Make sure that you specify a name if the resulting clone location as defined by
git would conflict as we see here. Both of these would have been cloned into
dotfiles which would cause the clone to fail for the second module if we didn't
specify name for chasinglogic's dotfiles.

An additional use for modules is that of a git repo you want to clone but not
link. An example use would be for downloading
[Spacemacs](https://github.com/syl20bnr/spacemacs) or any such community
configuration like oh-my-zsh, etc.

```yaml
modules:
    - repo: git@github.com:syl20bnr/spacemacs
      link: none
      pull_only: true
      location: ~/.emacs.d
```

Here we specify a few extra keys. There purpose should be self explanatory but
if you're curious [below](#available-keys) is a detailed explanation of all keys
that each module configuration supports.

#### Available Keys

- [repo](#repo)
- [name](#name)
- [location](#location)
- [link](#link)
- [pull\_only](#pull\_only)
- [mappings](#mappings)

##### repo

Required, this is the git repo to clone for the module.

##### name

This changes the cloned name. This only has an effect if location is not
provided. Normally a git repo would be cloned into
`$XDG_CONFIG_HOME/dfm/modules` and the resulting folder would be named whatever
git decides it should be based on the git url. If this is provided it will be
cloned into the modules directory with the specified name. This is useful if
multiple profiles use the same module.

##### location

If provided module will be cloned into the specified location. You can use the
`~` bash expansion here to represent `$HOME`. No other expansions are available.
This option is useful for cloning community configurations like oh-my-zsh or
spacemacs. 

##### link

Determines when to link the module. Link in this context means that it will be
treated like a normal dotfile profile, so all files will go through the same
translation rules as a regular profile and be linked accordingly. Available
values are "post", "pre", and "none". "post" is the default and means that the
module will be linked after the parent profile. "pre" means this will be linked
before the parent profile, use this if for instance you want to use most files
from this profile and override a few files with those from the parent file since
dfm will overwrite the links with the last one found. "none" means the module is
not a dotfiles profile and should not be linked at all, an example being
community configuration repositories like oh-my-zsh or spacemacs.

##### pull\_only

If set to `true` will not attempt to push any changes. Note that dfm always
tries to push to origin master, so if you don't have write access to the repo or
don't want it to automatically push to master then you should set this to true.
Useful for community configuration repositories.

##### mappings

A list of file mappings as described below in [Mappings](#mappings). Modules do
not inherit parent mappings, they do however inherit the default mappings as
described in [Skips Relevant Files](#skips-relevant-files)

### Mappings

Mappings are a way of defining custom file locations. For instance if dfm finds
a file or folder in your dotfile repo that is named `my_config` it will try and
"translate" that to `$HOME/.my_config` and create a symlink there pointing at
the one in your dotfile repo. With mappings you could replace this behavior and
make it so dfm will link .my_config wherever you wish. For example Visual Studio
Code has a non-standard (at least in the XDG\_CONFIG\_HOME / standard unix
dotfile sense) place to store it's configuration. So if you wanted to keep your
Visual Studio Code configuration in a dotfile repo you could add these mappings
to your `.dfm.yml` (assuming you're on a Mac):

```yaml
mappings:
  - match: my_vs_code_config
    is_dir: true
    dest: ~/Library/Application Support/Code/User/
  - match: .vscode
    skip: true
```

Here dfm would find my_\vs\_code\_config in your dotfile repo, recognize that it
has a mapping and perform the correct behavior based on your configuration. In
this case it would link each file inside of my\_vs\_code\_config and to the
corresponding file name in `~/Library/Application Support/Code/User`.
Additionally, it would skip any file that was an exact match for `.vscode` since
this directory we likely want in git but not linked via dfm.

#### Available Configuration

Filemaps support many configuration options:

- [match](#match)
- [skip](#skip)
- [dest](#dest)
- [regexp](#regexp)
- [is\_dir](#is\_dir)

##### match

Match is a required configuration option. This indicates what files match this
mapping. By default the file's actual name with no leading directories is taken
as a literal "is equal to" comparison. As an example if you have a file that's
in `$HOME/.config/dfm/profiles/my_profile/somefile` the correct match would be
simply `match: somefile` since dfm will do the comparison against the final file
name in the path. Note that this only works for files which are in the "root"
directory of your dotfile repo or in `$DOTFILE_ROOT/config`.

If more complicated match logic is required you can set [regexp](#regexp) to
true in which case the value of match is compiled into a regular expression
using the go standard library [regexp](https://golang.org/pkg/regexp/). See that
libraries docs for more info.

##### skip

If provided `dest` is ignored and the file is simply skipped if it matches
`match` instead of linked.

##### dest

Where to link the file instead of `$HOME`. The `~` expansion for `$HOME` is
supported here but no other expansions are available. This should be the
directory in which you want the file linked.

##### regexp

If true then `match` is treated as a 
[go regular expression](https://golang.org/pkg/regexp/)

##### is_dir

If is_dir is true then all files inside of the matched directory will be linked
at the new `dest` directory. If not provided and it is a directory that
directory would simply be linked using the normal logic. 

### Hooks

Hooks in dfm are used for those few extra tasks that you need to do whenever
your dotfiles are synced or cloned. One example is that I have ansible code in
my dotfiles that I like to run whenever I sync. So I wrote an `after_sync` hook
to do just that:

```yaml
hooks:
  after_sync:
    - ansible-playbook ~/.config/dfm/profiles/chasinglogic/ansible/dev-mac.yml
```

Now whenever I sync my dotfiles ansible will run my dev-mac playbook to make
sure that my packages etc are also in sync!

You can wrap any dfm subcommand. `hooks` is simply map where the key is when to
run the hook. Which is specified using the syntax `before_{command_name}` or
`after_{command_name}`. The value of the key should then be a yaml list of bash
commands to run. A slightly more complex example (though a completely useless
one):

```yaml
hooks:
  after_link:
    - ls -l
    - whoami
    - echo "All done!"
```

The current process environment is passed down to the shell so you can use
`$HOME` etc environment variables in your commands.

**NOTE:** the before and after clone hooks rarely work how you expect.
`after_clone` will be run from the freshly cloned repo. However if there is
already a current profile then both `before_clone` and `after_clone` will be run
from whatever the current profile is. It's best to avoid writing hooks around
clone since this behavior is hard to intuit. 

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
