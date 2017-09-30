# dfm
A dotfile manager for lazy people and pair programmers.

Dotfile Manager will allow you to create "profiles" for dotfiles underneath one
unix account and easily switch between them. It requires a git repo with your
dotfiles in it and that the dotfiles be placed how you want them represented in
your home directory.

It does not require that the dotfiles in your repo start with dots though it
handles either case equally well.

## Installation
The easiest (and currently only) way to install dfm is using go

```bash
$ go get github.com/chasinglogic/dfm
```

## Updating

Make sure dfm is updated to bring in the latest bug fixes and features. If you are installing dfm for the first time, you can skip this step.

```bash
go get -u github.com/chasinglogic/dfm
```

## Usage

```text
Dotfile management written for pair programmers. Examples on getting
started with dfm are avialable at https://github.com/chasinglogic/dfm

Usage:
  dfm [flags]
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
  where       prints the current profile directory path

Flags:
  -d, --dry-run   don't make changes just print what would happen
  -h, --help      help for dfm
  -v, --verbose   verbose output

Use "dfm [command] --help" for more information about a command.
```

dfm is mostly a thin wrapper around git and just manages repos and symlinks
for you. As such most of the dfm commands are directly analogous to git
commands.

> **A note about the $XDG\_CONFIG\_HOME (commonly $HOME/.config) directory:**
>
> dfm respects dotfiles which exist in the $XDG\_CONFIG\_HOME directory, meaning
> if in your repo you have a folder named config or .config it will translate
> those into the  $XDG\_CONFIG\_HOME directory appropriately. Similarly when
> using `dfm add` if inside your $XDG\_CONFIG\_HOME or $HOME/.config directories
> it will add those to the repo appropriately.

### Quick Start (Existing dotfiles repo with the Git backend)

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

Once you have made some changes you can simply sync your changes to git with:

```bash
dfm sync
```

See the Usage Notes below for some quick info on what to expect from other dfm
commands.

### Quick Start (No existing dotfiles repo)

If you do not have a dotfiles repo the best place to start is with `dfm init`

```bash
dfm init my-new-profile
```

Then run `dfm link` to set it as the active profile, this is also how you
switch profiles

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

Then set the remote to your remote repo

**Note:** When creating the remote repo do not choose any options such as
"initialize this repo with a README" otherwise git will get cranky when you add
the remote because of a recent git update and how it handles [unrelated
histories](http://stackoverflow.com/questions/37937984/git-refusing-to-merge-unrelated-histories)
if you do don't worry the linked post explains how to get past it.

```bash
dfm git remote add origin https://github.com/myusername/dotfiles
```

Then simply sync your changes

```bash
dfm sync
```

Now you're done!

### A quick note about git commands and flags

dfm does always simply push your commands directly through to git.

The git sub command will push directly through so you can run whatever you want
as if you were in that directory. For example:

`dfm git checkout -b my-work-laptop`

Stdin, Stdout, and Stderr are given to git so even things that run a text
editor for example will work as expected.

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
