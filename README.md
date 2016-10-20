# dfm
A dotfile manager for lazy people and pair programmers.

Dotfile Manager will allow you to create "profiles" for dotfiles underneath one
unix account and easily switch between them. It requires a git repo with your
dotfiles in it and that the dotfiles be placed how you want them represented in
your home directory.

It does not require that the dotfiles in your repo start with dots though it 
handles either case equally well.

## Installation
The easiest (and currently only) way to install dfm is using the go tool

```bash
go get github.com/chasinglogic/dfm/cmd/dfm
```

## Usage

```
AME:
   dfm - Manage dotfiles.

USAGE:
   dfm [global options] command [command options] [arguments...]

VERSION:
   1.0-dev

AUTHOR(S):
   Mathew Robinson <mathew.robinson3114@gmail.com>

COMMANDS:
     add, a      Add a file to the current profile.
     clone, c    Create a dotfiles profile from a git repo.
     link, l     Recreate the links from the dotfiles profile.
     list, ls    List available profiles
     pull, pl    Pull the latest version of the profile from origin master.
     push, ps    Push your local version of the profile to the remote.
     remove, rm  Remove the profile and all it's symlinks.
     init, i     Create a new profile with `NAME`
     commit, cm  Runs git commit for the profile using `MSG` as the message
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config DIR, -c DIR  Use DIR for storing dfm configuration and profiles (default: "/Users/mr91060/.config/dfm")
   --verbose, --vv       Print verbose messaging.
   --dry-run             Don't create symlinks just print what would be done.
   --help, -h            show help
   --version, -v         print the version
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

`
