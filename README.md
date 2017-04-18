# dfm
A dotfile manager for lazy people and pair programmers.

Dotfile Manager will allow you to create "profiles" for dotfiles underneath one
unix account and easily switch between them. It requires a git repo with your
dotfiles in it and that the dotfiles be placed how you want them represented in
your home directory.

It does not require that the dotfiles in your repo start with dots though it 
handles either case equally well.

## Installation
The easiest (and currently only) way to install dfm is using pip

```bash
$ pip3 install dfm
```

To install in the global path you'll need to run `pip3` as root.

```
# pip3 install dfm
```

## Usage

```text
Usage: dfm [OPTIONS] COMMAND [ARGS]...

  A dotfile manager for lazy people and pair programmers.

Options:
  -vv, --verbose
  -c, --config PATH  The path where dfm stores it's config and profiles.
  --help             Show this message and exit.

Commands:
  add       Add a file or directory to the current...
  checkout  Switch to a different branch for the active...
  clone     Clone a profile from a git repo.
  commit    Run a git commit for the current profile.
  git       Run the given git command in the current...
  init      Create an empty profile with the given name.
  license   Show dfm licensing info.
  link      Link the profile with the given name.
  pull      Pull changes from the remote.
  push      Push local changes to the remote.
  remote    Set the git remote for the current profile.
  rm        Remove the profile with the given name.
  upgrade   Upgrade from the old style config to the new...
  version   Show the current dfm version.
  where     Return the path to the current profile.
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
dfm commit -am "init dotfile repo"
```

Then set the remote to your remote repo

**Note:** When creating the remote repo do not choose any options such as
"initialize this repo with a README" otherwise git will get cranky when you add
the remote because of a recent git update and how it handles [unrelated
histories](http://stackoverflow.com/questions/37937984/git-refusing-to-merge-unrelated-histories)
if you do don't worry the linked post explains how to get past it.

```bash
dfm remote https://github.com/myusername/dotfiles
```

Then simply push them up

```bash
dfm push
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

This code is distributed under the Apache 2.0 License.

```
Copyright 2016 Mathew Robinson

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
