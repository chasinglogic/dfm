mod profiles;

use std::{
    env,
    ffi::OsString,
    fs::{self, File},
    io::{self, BufReader},
    path::{Path, PathBuf},
    process,
};

use clap::{command, crate_version, Parser, Subcommand};
use profiles::profile::Profile;
use walkdir::WalkDir;

#[derive(Debug, Parser)]
#[command(
    name = "dfm",
    about = "A dotfile manager for pair programmers and lazy people.", 
    long_about = None,
    version = crate_version!(),
)]
struct CLI {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Debug, Subcommand)]
enum Commands {
    Test,
    #[command(visible_alias = "w")]
    Where,
    #[command(visible_alias = "st")]
    Status,
    #[command(visible_alias = "g")]
    Git {
        #[arg(trailing_var_arg = true, allow_hyphen_values = true)]
        args: Vec<String>,
    },
    #[command(visible_alias = "ls")]
    List,
    #[command(visible_alias = "l")]
    Link {
        // New profile to switch to and link
        #[arg(default_value_t)]
        profile_name: String,
    },
    #[command(visible_alias = "i")]
    Init {
        #[arg(required = true)]
        profile_name: String,
    },
    #[command(visible_alias = "rm")]
    Remove {
        #[arg(required = true)]
        profile_name: String,
    },
    #[command(visible_alias = "rh")]
    RunHook {
        #[arg(required = true)]
        hook_name: String,
    },
    // TODO:
    // Clone
    // Add
    // Sync
    // Clean
    #[command(external_subcommand)]
    External(Vec<OsString>),
}

#[derive(Debug, serde::Deserialize, serde::Serialize)]
struct State {
    current_profile: String,
}

impl Default for State {
    fn default() -> Self {
        State {
            current_profile: "".to_string(),
        }
    }
}

impl State {
    fn load(fp: &Path) -> Result<State, io::Error> {
        let fh = File::open(fp)?;
        let buffer = BufReader::new(fh);
        Ok(serde_json::from_reader(buffer)?)
    }

    fn save(&self, filepath: &Path) -> Result<(), io::Error> {
        if let Some(parent) = filepath.parent() {
            if !parent.exists() {
                fs::create_dir_all(parent).expect("Unable to create dfm directory!");
            }
        }

        let file_handle = File::create(filepath)?;
        Ok(serde_json::to_writer(file_handle, self)?)
    }
}

fn dfm_dir() -> PathBuf {
    let home = env::var("HOME").unwrap_or("".to_string());
    let mut path = PathBuf::from(home);
    path.push(".config");
    path.push("dfm");
    path
}

fn state_file() -> PathBuf {
    let mut state_fp = dfm_dir();
    state_fp.push("state.json");
    state_fp
}

fn profiles_dir() -> PathBuf {
    let mut path = dfm_dir();
    path.push("profiles");
    path
}

fn load_profile(name: &str) -> Profile {
    let mut path = profiles_dir();
    path.push(name);
    Profile::load(&path)
}

fn force_available(profile: Option<Profile>) -> Profile {
    match profile {
        None => {
            eprintln!("No profile is currently loaded!");
            process::exit(1);
        }
        Some(p) => p,
    }
}

fn main() {
    let args = CLI::parse();
    let state_fp = state_file();
    let mut state = match State::load(&state_fp) {
        Ok(state) => state,
        Err(err) => match err.kind() {
            io::ErrorKind::NotFound => State::default(),
            _ => panic!("{}", err),
        },
    };
    let current_profile: Option<Profile> = if state.current_profile != "" {
        Some(load_profile(&state.current_profile))
    } else {
        None
    };

    match args.command {
        Commands::Test => match current_profile {
            Some(profile) => println!("{:#?}", profile.name()),
            None => println!("Current profile not loaded!"),
        },
        Commands::Where => println!("{}", force_available(current_profile).config.location),
        Commands::List => {
            for entry in WalkDir::new(profiles_dir()).min_depth(1).max_depth(1) {
                println!("{}", entry.unwrap().file_name().to_string_lossy());
            }
        }
        Commands::Git { args } => force_available(current_profile)
            .git(args)
            .map(|_| ())
            .expect("Unable to run git on the current profile!"),
        Commands::RunHook { hook_name } => force_available(current_profile)
            .run_hook(&hook_name)
            .expect("Unable to run hook!"),
        Commands::Link { profile_name } => {
            let new_profile = if profile_name != "" {
                load_profile(&profile_name)
            } else {
                force_available(current_profile)
            };
            new_profile.link().expect("Error linking profile!");
            state.current_profile = new_profile.name();
        }
        Commands::Init { profile_name } => {
            let mut path = profiles_dir();
            path.push(&profile_name);
            if path.exists() {
                eprintln!(
                    "Unable to create profile as {} already exists!",
                    path.to_string_lossy()
                );
                process::exit(1);
            }

            fs::create_dir_all(&path).expect("Unable to create profile directory!");
            let new_profile = Profile::load(&path);
            new_profile.init().expect("Error initialising profile!");
        }
        Commands::Remove { profile_name } => {
            let mut path = profiles_dir();
            path.push(&profile_name);
            if !path.exists() {
                eprintln!("No profile with exists at path: {}", path.to_string_lossy());
                process::exit(1);
            }

            if !path.is_dir() {
                eprintln!("Profile exists but is not a directory!");
                process::exit(1);
            }

            fs::remove_dir_all(&path).expect("Unable to remove profile directory!");
            println!("Profile {} successfully removed.", profile_name);
        }
        Commands::Status => force_available(current_profile)
            .status()
            .map(|_| ())
            .expect("Unexpected error running git!"),
        Commands::External(args) => {
            let plugin_name = format!("dfm-{}", args[0].to_str().unwrap_or_default());
            println!("Calling out to {:?} with {:?}", plugin_name, &args[1..]);
        }
    }

    state.save(&state_fp).expect("Unable to save state!")
}
