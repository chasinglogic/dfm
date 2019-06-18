extern crate clap;

use std::env;
use std::error::Error;
use std::fs::File;
use std::io;
use std::io::{Read, Write};
use std::path::PathBuf;
use std::process;

use clap::{App, AppSettings, Arg, SubCommand};
use serde::{Deserialize, Serialize};

use crate as dfm;

const VERSION: &str = env!("CARGO_PKG_VERSION");

fn xdg_dir() -> PathBuf {
    match env::var("XDG_CONFIG_HOME") {
        Ok(path) => PathBuf::from(path),
        Err(_) => {
            let home = env::var("HOME").unwrap_or("".to_string());
            let mut home_p = PathBuf::from(home);
            home_p.push(".config");
            home_p
        }
    }
}

fn dfm_dir() -> PathBuf {
    match env::var("DFM_CONFIG_DIR") {
        Ok(path) => PathBuf::from(path),
        Err(_) => {
            let mut path = xdg_dir();
            path.push("dfm");
            path
        }
    }
}

fn state_file_p() -> PathBuf {
    let mut p = dfm_dir();
    p.push("state.json");
    p
}

fn profile_storage_dir() -> PathBuf {
    let mut p = dfm_dir();
    p.push("profiles");
    p
}

fn profile_dir(name: &str) -> PathBuf {
    let mut storage = profile_storage_dir();
    storage.push(name);
    storage
}

#[derive(Deserialize, Serialize)]
struct State {
    current_profile: String,
}

impl State {
    fn load() -> Option<State> {
        let sf = state_file_p();
        let mut contents = String::new();
        match File::open(sf) {
            Ok(mut f) => {
                if let Err(_) = f.read_to_string(&mut contents) {
                    return None;
                };

                match serde_yaml::from_str(&contents) {
                    Ok(s) => Some(s),
                    Err(_) => None,
                }
            }
            Err(_) => None,
        }
    }

    fn save(self) -> Result<(), io::Error> {
        let sf = state_file_p();
        let display = sf.display();

        let mut file = match File::create(&sf) {
            Err(why) => {
                println!(
                    "Couldn't save state file {}: {}",
                    display,
                    why.description()
                );
                process::exit(1);
            }
            Ok(file) => file,
        };

        let bytes = serde_yaml::to_vec(&self).expect("Failed to serialize app state");
        file.write_all(&bytes)
    }
}

fn main() {
    let matches = App::new("dfm")
        .version(VERSION)
        .author("Mathew Robinson <chasinglogic@gmail.com>")
        .about("A dotfile manager for pair programmers and lazy people")
        .help(
            "Examples on getting started with dfm are avialable at https://github.com/chasinglogic/dfm",
        )
        .subcommand(
            SubCommand::with_name("sync")
                .alias("s")
                .about("Sync dotfiles and modules")
        )
        .subcommand(
            SubCommand::with_name("add")
                .alias("a")
                .about("Add a file to the current dotfile profile")
                .arg(
                    Arg::with_name("file")
                        .multiple(true)
                )
        )
        .subcommand(
            SubCommand::with_name("clean")
                .alias("x")
                .about("Clean dead symlinks")
        )
        .subcommand(
            SubCommand::with_name("git")
                .alias("g")
                .about("Run the given git command on the current profile")
                .setting(AppSettings::TrailingVarArg)
                .arg(Arg::with_name("cmd").multiple(true))
        )
        .subcommand(
            SubCommand::with_name("clone")
                .alias("c")
                .about("Use git clone to download a dotfile profile")
        )
        .subcommand(
            SubCommand::with_name("init")
                .alias("i")
                .about("Create a new profile")
        )
        .subcommand(
            SubCommand::with_name("link")
                .alias("l")
                .about("Link and activate a dotfile profile")
        )
        .subcommand(
            SubCommand::with_name("list")
                .alias("ls")
                .about("List available profiles")
        )
        .subcommand(
            SubCommand::with_name("remove")
                .alias("rm")
                .about("Remove a profile")
        )
        .subcommand(
            SubCommand::with_name("run-hook")
                .alias("rh")
                .about("Run dfm hooks without runnign the associated command")
        )
        .subcommand(
            SubCommand::with_name("where")
                .alias("w")
                .about("Prints the location of the active dotfile profile")
        )
        .get_matches();
}
