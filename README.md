# **D**otfile **F**ile **M**anager
A dotfile manager for lazy people and pair programmers.

Dot File Manager will allow you to create "profiles" for dotfiles underneath one
unix account and easily switch between them. It requires a git repo with your
dotfiles in it and that the dotfiles be placed how you want them represented in
your home directory.

It does not require that the dotfiles in your repo start with dots though it 
handles either case equally well.

## Installation
The easiest (and currently only) way to install dfm is using pip

```bash
pip install dfm
```

If you're on ubuntu or a system that ships with python 2 you probably need
to use pip3

```bash
pip3 install dfm
```

## Usage

```
Usage: dfm [OPTIONS] COMMAND [ARGS]...

  A dotfile manager for lazy people and pair programmers.

Options:
  -vv, --verbose
  -c, --config PATH  The path where dfm stores it's config and profiles.
  --help             Show this message and exit.

Commands:
  add      Add a file or directory to the current...
  chk      Switch to a different branch for the active...
  clone    Clone a profile from a git repo.
  commit   Run a git commit for the current profile.
  init     Create an empty profile with the given name.
  license  Show dfm licensing info.
  link     Link the profile with the given name.
  pull     Pull changes from the remote.
  push     Push local changes to the remote.
  rm       Remove the profile with the given name.
  version  Show the current dfm version.
```

dfm is mostly a thin wrapper around git and just manages repos and symlinks 
for you. As such most of the dfm commands are directly analogous to git 
commands.

A note about the $XDG\_CONFIG\_HOME (commonly $HOME/.config) directory:

dfm respects dotfiles which exist in the $XDG\_CONFIG\_HOME directory, meaning 
if in your repo you have a folder named config or .config it will translate 
those into the  $XDG\_CONFIG\_HOME directory appropriately. Similarly when 
using `dfm add` if inside your $XDG\_CONFIG\_HOME or $HOME/.config directories 
it will add those to the repo appropriately.

### Quick Start (Existing dotfiles repo)

If you already have a dotfiles repo you can start by cloning it using the clone
command. **Note:** ssh urls work as well since it's just passed to git.

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

Then simply run `dfm commit` to commit them

```bash
dfm commit "init dotfile repo"
```

Then set the remote to your remote repo

```bash
dfm remote https://github.com/myusername/dotfiles
```

Then simply push them up

```bash
dfm push
```

Now you're done!

### A quick note about git commands and flags

dfm does not simply push your commands directly through to git, I wanted to
avoid that for those who are not or do not want to become familiar with git
I've thought about adding this feature later for more advanced users but for
now I think this is fine.

Another note is that when running `dfm commit` we give git the -a flag so all 
of your changes will be added to that commit. 

## Contributing

1. Fork it! :fork_and_knife:
2. Create an issue describing what you're working on.
3. Create your feature branch: `git checkout -b my-new-feature`
4. Commit your changes: `git commit -am 'Add some feature'`
5. Push to the branch: `git push origin my-new-feature`
6. :fire: Submit a pull request :D :fire:

All pull requests should go to the develop branch not master. Thanks!

## License

dfm is distributed under the GPLv3

```
  dfm, a dotfile manager for lazy people and pair programmers

  Copyright (C) 2016 Mathew Robinson <mathew.robinson3114@gmail.com>

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
