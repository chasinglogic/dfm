mod profiles;

use shellexpand;
use std::{
    default, env,
    ffi::OsString,
    fs::File,
    io::{self, BufReader},
    path::{Path, PathBuf},
};

use self::profiles::profile::Profile;
use clap::{command, Parser, Subcommand};

#[derive(Debug, Parser)]
#[command(name = "dfm")]
#[command(about = "A dotfile manager for pair programmers and lazy people.", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Debug, Subcommand)]
enum Commands {
    #[command()]
    Test,
    #[command()]
    Where,
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
        // TODO: Create intermediate directories if required.
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
    let home = env::var("HOME").unwrap_or("".to_string());
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

fn main() {
    let args = Cli::parse();
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
            Some(profile) => println!("{:#?}", profile),
            None => println!("Current profile not loaded!"),
        },
        Commands::Where => match current_profile {
            None => eprintln!("No profile is currently selected!"),
            Some(profile) => println!("{}", profile.config.location),
        },
        Commands::External(args) => {
            let plugin_name = format!("dfm-{}", args[0].to_str().unwrap_or_default());
            println!("Calling out to {:?} with {:?}", plugin_name, &args[1..]);
        }
    }

    state.save(&state_fp).expect("Unable to save state!")
}
