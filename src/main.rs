extern crate clap;
extern crate regex;
extern crate serde;
extern crate serde_yaml;
extern crate walkdir;
#[macro_use]
extern crate text_io;

use std::fs;
use std::process;

use clap::{App, AppSettings, Arg, SubCommand};

mod hooks;
mod link;
mod mapping;
mod profile;
mod repo;
mod state;
mod util;

const VERSION: &str = env!("CARGO_PKG_VERSION");

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
                .arg(Arg::with_name("profile"))
        )
        .subcommand(
            SubCommand::with_name("add")
                .alias("a")
                .about("Add a file to the current dotfile profile")
                .arg(Arg::with_name("profile").long("profile").short("p"))
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
                .arg(Arg::with_name("profile").long("profile").short("p"))
                .arg(Arg::with_name("cmd").multiple(true))
        )
        .subcommand(
            SubCommand::with_name("clone")
                .alias("c")
                .arg(Arg::with_name("name").long("name").short("n"))
                .about("Use git clone to download a dotfile profile")
        )
        .subcommand(
            SubCommand::with_name("init")
                .alias("i")
                .about("Create a new profile")
                .arg(Arg::with_name("profile"))
        )
        .subcommand(
            SubCommand::with_name("link")
                .alias("l")
                .about("Link and activate a dotfile profile")
                .arg(Arg::with_name("profile"))
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
                .arg(
                    Arg::with_name("profile")
                        .multiple(true)
                )
        )
        .subcommand(
            SubCommand::with_name("run-hook")
                .alias("rh")
                .about("Run dfm hooks without runnign the associated command")
                .arg(Arg::with_name("profile").long("profile").short("p"))
        )
        .subcommand(
            SubCommand::with_name("where")
                .alias("w")
                .about("Prints the location of the active dotfile profile")
                .arg(Arg::with_name("profile").long("profile").short("p"))
        )
        .get_matches();

    let mut state = match state::State::load() {
        Some(s) => s,
        None => state::State::default(),
    };

    match matches.subcommand() {
        ("add", Some(args)) => unimplemented!(),
        ("clone", Some(args)) => {
            let repo = args.value_of("repo").unwrap();
            let name = args.value_of("name").unwrap();
            let mut child = process::Command::new("git");
            child.stdin(process::Stdio::inherit());
            child.stdout(process::Stdio::inherit());
            child.stderr(process::Stdio::inherit());
            child.args(&["clone", repo, name]);
            let mut proc = child.spawn().expect("Unable to run git clone");
            proc.wait();
        }
        ("git", Some(args)) => {
            let profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or(state.current_profile.as_str()),
            );

            let cmd: Vec<&str> = args
                .values_of("cmd")
                .unwrap_or(clap::Values::default())
                .collect();
            match profile.repo.git(&cmd) {
                Ok(_) => (),
                Err(e) => {
                    println!("Unexpected error: {}", e);
                    process::exit(1);
                }
            }
        }
        ("list", Some(args)) => {
            let profiles_dir = util::profile_storage_dir();
            for res in fs::read_dir(&profiles_dir).expect("Unable to read profile directory") {
                if let Ok(entry) = res {
                    println!("{}", entry.file_name().to_string_lossy());
                }
            }
        }
        ("link", Some(args)) => {
            let name = args
                .value_of("profile")
                .unwrap_or(state.current_profile.as_str());
            let profile = util::load_profile(name);
            profile.link();
            state.current_profile = name.to_string();
        }
        ("run-hook", Some(args)) => {
            let profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or(state.current_profile.as_str()),
            );
            profile
                .hooks
                .run(&profile.repo.path, args.value_of("hook").unwrap_or(""));
        }
        ("sync", Some(args)) => {
            let profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or(state.current_profile.as_str()),
            );
            profile.sync();
        }
        ("where", Some(args)) => {
            let profile_dir = util::profile_dir(&state.current_profile);
            println!("{}", profile_dir.display());
        }
        ("remove", Some(args)) => {
            let dir = util::profile_dir(
                args.value_of("profile")
                    .unwrap_or(state.current_profile.as_str()),
            );
            fs::remove_dir_all(&dir);
        }
        (s, _) => {
            println!("Not a valid subcommand: {}", s);
            println!("{}", matches.usage());
            process::exit(1);
        }
    }

    state.save().expect("Unable to write new state file");
}
