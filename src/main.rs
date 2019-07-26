extern crate clap;
extern crate regex;
extern crate serde;
extern crate serde_yaml;
extern crate walkdir;
#[macro_use]
extern crate text_io;
extern crate log;
extern crate simplelog;

use std::fs;
use std::path::Path;
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
        .about(
            "A dotfile manager for pair programmers and lazy people

Examples on getting started with dfm are avialable at https://github.com/chasinglogic/dfm",
        )
        .bin_name("dfm")
        .settings(&[
            AppSettings::GlobalVersion,
            AppSettings::SubcommandRequiredElseHelp,
        ])
        .arg(
            Arg::with_name("config-dir")
                .long("config-dir")
                .short("c")
                .takes_value(true),
        )
        .arg(Arg::with_name("verbose").long("verbose").takes_value(false))
        .subcommand(
            SubCommand::with_name("sync")
                .alias("s")
                .about("Sync dotfiles and modules")
                .arg(
                    Arg::with_name("profile")
                        .short("p")
                        .long("profile")
                        .help("The profile to operate on")
                        .takes_value(true),
                )
                .arg(
                    Arg::with_name("message")
                        .long("message")
                        .short("m")
                        .takes_value(true),
                ),
        )
        .subcommand(
            SubCommand::with_name("add")
                .alias("a")
                .about("Add a file to the current dotfile profile")
                .arg(
                    Arg::with_name("profile")
                        .long("profile")
                        .short("p")
                        .help("The profile to operate on"),
                )
                .arg(Arg::with_name("file").multiple(true)),
        )
        .subcommand(
            SubCommand::with_name("git")
                .alias("g")
                .about("Run the given git command on the current profile")
                .setting(AppSettings::TrailingVarArg)
                .arg(
                    Arg::with_name("profile")
                        .long("profile")
                        .short("p")
                        .help("The profile to operate on"),
                )
                .arg(Arg::with_name("cmd").multiple(true)),
        )
        .subcommand(
            SubCommand::with_name("clone")
                .alias("c")
                .arg(
                    Arg::with_name("name")
                        .long("name")
                        .short("n")
                        .required(true)
                        .takes_value(true),
                )
                .arg(
                    Arg::with_name("url")
                        .takes_value(true)
                        .required(true)
                        .takes_value(true),
                )
                .about("Use git clone to download a dotfile profile"),
        )
        .subcommand(
            SubCommand::with_name("init")
                .alias("i")
                .about("Create a new profile")
                .arg(
                    Arg::with_name("name")
                        .takes_value(true)
                        .required(true)
                        .help("Name of the new profile"),
                ),
        )
        .subcommand(
            SubCommand::with_name("link")
                .alias("l")
                .about("Link and activate a dotfile profile")
                .arg(Arg::with_name("overwrite").short("o").long("overwrite"))
                .arg(Arg::with_name("profile")),
        )
        .subcommand(
            SubCommand::with_name("list")
                .alias("ls")
                .about("List available profiles"),
        )
        .subcommand(
            SubCommand::with_name("remove")
                .alias("rm")
                .about("Remove a profile")
                .arg(Arg::with_name("profile").multiple(true)),
        )
        .subcommand(
            SubCommand::with_name("run-hook")
                .alias("rh")
                .about("Run dfm hooks without runnign the associated command")
                .arg(Arg::with_name("profile").long("profile").short("p")),
        )
        .subcommand(
            SubCommand::with_name("where")
                .alias("w")
                .about("Prints the location of the active dotfile profile")
                .arg(Arg::with_name("profile").long("profile").short("p")),
        )
        .get_matches();

    let cfd = util::cfg_dir(matches.value_of("config-dir").map(Path::new));
    let sf = util::state_file_p(&cfd);
    let mut state = match state::State::load(&sf) {
        Some(s) => s,
        None => state::State::default(),
    };

    if matches.is_present("verbose") {
        simplelog::SimpleLogger::init(simplelog::LevelFilter::Debug, simplelog::Config::default())
            .unwrap();
    }

    match matches.subcommand() {
        ("add", Some(args)) => {
            let profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );
            let profile_dir = util::profile_dir(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );

            let files = args.values_of("file").unwrap();
            for file in files {
                let fp = Path::new(file);
                let new_file = match fp.strip_prefix(&profile.target_dir) {
                    Ok(f) => f,
                    Err(_) => {
                        println!("file does not belong to profile's target dir, cannot add automatically");
                        process::exit(1);
                    }
                };
                let mut pd = profile_dir.clone();
                pd.push(new_file);
                if let Err(e) = fs::copy(fp, &pd) {
                    println!("unable to copy {} to {}: {}", fp.display(), pd.display(), e);
                    process::exit(1);
                }

                if let Err(e) = if fp.is_dir() {
                    fs::remove_dir_all(fp)
                } else {
                    fs::remove_file(fp)
                } {
                    println!("unable to remove {}: {}", fp.display(), e);
                    process::exit(1);
                }
            }

            if let Err(e) = profile.link(false) {
                println!(
                    "Error linking profile {}: {}",
                    profile.repo.path.display(),
                    e
                );
                process::exit(1);
            }
        }
        ("init", Some(args)) => {
            let profile_dir = util::profile_dir(args.value_of("name").unwrap(), &cfd);
            println!("Creating {}", profile_dir.display());
            util::ensure_exists(&profile_dir);
            let profile = util::load_profile(args.value_of("name").unwrap(), &cfd);
            if let Err(e) = profile.init() {
                println!("unable to initialize profile: {}", e);
                process::exit(1);
            }
        }
        ("clone", Some(args)) => {
            let repo = args.value_of("url").unwrap();
            let name = args.value_of("name").unwrap();
            let new_profile_dir = util::profile_dir(name, &cfd);
            let mut child = process::Command::new("git");
            child.stdin(process::Stdio::inherit());
            child.stdout(process::Stdio::inherit());
            child.stderr(process::Stdio::inherit());
            child.args(&["clone", repo, &new_profile_dir.to_string_lossy()]);
            let mut proc = child.spawn().expect("Unable to run git clone");
            if proc.wait().is_err() {
                process::exit(1)
            };
        }
        ("git", Some(args)) => {
            let profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );

            let cmd: Vec<&str> = args.values_of("cmd").unwrap_or_default().collect();
            if let Err(e) = profile.repo.git(&cmd) {
                println!("Unexpected error: {}", e);
                process::exit(1);
            }
        }
        ("list", Some(_)) => {
            let profiles_dir = util::profile_storage_dir(&cfd);
            for res in fs::read_dir(&profiles_dir).expect("Unable to read profile directory") {
                if let Ok(entry) = res {
                    println!("{}", entry.file_name().to_string_lossy());
                }
            }
        }
        ("link", Some(args)) => {
            let name = args
                .value_of("profile")
                .unwrap_or_else(|| state.current_profile.as_str());
            let profile = util::load_profile(name, &cfd);

            if let Err(e) = profile.link(args.is_present("overwrite")) {
                println!(
                    "Error linking profile {}: {}",
                    profile.repo.path.display(),
                    e
                );
                process::exit(1);
            }

            state.current_profile = name.to_string();
        }
        ("run-hook", Some(args)) => {
            let profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );

            let hook = args.value_of("hook").unwrap_or("");
            if let Err(e) = profile.hooks.run(&profile.repo.path, hook) {
                println!(
                    "Error running hook {} in {}: {}",
                    hook,
                    profile.repo.path.display(),
                    e
                );
                process::exit(1);
            }
        }
        ("sync", Some(args)) => {
            let mut profile = util::load_profile(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );

            if let Some(msg) = args.value_of("message") {
                profile.set_commit_msg(msg);
                profile.set_prompt(false);
            };

            if let Err(e) = profile.sync() {
                println!(
                    "Error syncing profile {}: {}",
                    profile.repo.path.display(),
                    e
                );
                process::exit(1);
            };
        }
        ("where", Some(args)) => {
            let profile_dir = util::profile_dir(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );
            println!("{}", profile_dir.display());
        }
        ("remove", Some(args)) => {
            let dir = util::profile_dir(
                args.value_of("profile")
                    .unwrap_or_else(|| state.current_profile.as_str()),
                &cfd,
            );

            if let Err(e) = fs::remove_dir_all(&dir) {
                println!("Unable to remove profile {}: {}", dir.display(), e);
                process::exit(1);
            }
        }
        (s, _) => {
            println!("Not a valid subcommand: {}", s);
            println!("{}", matches.usage());
            process::exit(1);
        }
    }

    state.save(&sf).expect("Unable to write new state file");
}
