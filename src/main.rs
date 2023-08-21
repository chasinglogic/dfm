mod cli;
mod profiles;

use crate::cli::state::{force_available, profiles_dir};

use std::{
    env,
    ffi::OsString,
    fs, io,
    path::Path,
    process::{self, Command},
};

use clap::{command, crate_version, CommandFactory, Parser, Subcommand, ValueEnum};
use clap_complete::{generate, Shell};
use profiles::Profile;
use walkdir::WalkDir;

#[derive(Debug, Parser)]
#[command(
    name = "dfm",
    about = "A dotfile manager for pair programmers and lazy people.", 
    long_about = "Dotfile management written for pair programmers. 
Examples on getting started with dfm are available at https://github.com/chasinglogic/dfm",
    version = crate_version!(),
)]
struct CLI {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Debug, Subcommand)]
enum Commands {
    #[command(
        visible_alias = "w",
        about = "Prints the location of the current dotfile profile"
    )]
    Where,
    #[command(
        visible_alias = "st",
        about = "Print the git status of the current dotfile profile"
    )]
    Status,
    #[command(
        visible_alias = "g",
        about = "Run the given git command on the current profile"
    )]
    Git {
        #[arg(trailing_var_arg = true, allow_hyphen_values = true)]
        args: Vec<String>,
    },
    #[command(
        visible_alias = "ls",
        about = "List available dotfile profiles on this system"
    )]
    List,
    #[command(
        visible_alias = "l",
        about = "Create links for a profile",
        long_about = "Creates symlinks in HOME for a dotfile Profile and makes it the active profile"
    )]
    Link {
        // New profile to switch to and link
        #[arg(
            default_value_t,
            help = "The profile name to link, if none given relinks the current profile"
        )]
        profile_name: String,
        #[arg(
            default_value_t,
            short,
            long,
            help = "If provided dfm will delete files and directories which exist at the target \
                    link locations. DO NOT USE THIS IF YOU ARE UNSURE AS IT WILL RESULT IN DATA LOSS"
        )]
        overwrite: bool,
    },
    #[command(visible_alias = "i", about = "Create a new profile")]
    Init {
        #[arg(required = true, help = "Name of the profile to create")]
        profile_name: String,
    },
    #[command(visible_alias = "rm", about = "Remove a profile")]
    Remove {
        #[arg(required = true, help = "Name of the profile to remove")]
        profile_name: String,
    },
    #[command(
        visible_alias = "rh",
        about = "Run dfm hooks without using normal commands",
        long_about = "Runs a hook without the need to invoke the side effects of a dfm command"
    )]
    RunHook {
        #[arg(required = true)]
        hook_name: String,
    },
    #[command(visible_alias = "s", about = "Sync your dotfiles")]
    Sync {
        #[arg(
            default_value_t,
            short,
            long,
            help = "Use the given message as the commit message"
        )]
        message: String,
    },
    #[command(about = "Use git clone to download an existing profile")]
    Clone {
        #[arg(required = true)]
        url: String,
        #[arg(
            default_value_t,
            short,
            long,
            help = "Name of the profile to create, defaults to the basename of <url>"
        )]
        name: String,
        #[arg(
            default_value_t,
            short,
            long,
            help = "If provided the profile will be immediately linked"
        )]
        link: bool,
        #[arg(
            default_value_t,
            short,
            long,
            help = "If provided dfm will delete files and directories which exist at the target \
                    link locations. DO NOT USE THIS IF YOU ARE UNSURE AS IT WILL RESULT IN DATA LOSS"
        )]
        overwrite: bool,
    },
    #[command(about = "Clean dead symlinks. Will ignore symlinks unrelated to DFM.")]
    Clean,
    #[command(
        about = "Add files to the current dotfile profile",
        long_about = "Add files to the current dotfile profile"
    )]
    Add {
        #[arg(required = true)]
        files: Vec<String>,
    },
    #[command(about = "Generate shell completions and print them to stdout")]
    GenCompletions {
        #[arg(required = true, help = "The shell to generate completions for.")]
        shell: String,
    },
    #[command(external_subcommand)]
    External(Vec<OsString>),
}

fn main() {
    let args = CLI::parse();
    let mut state = cli::state::load_or_default();

    let current_profile: Option<Profile> = if state.current_profile != "" {
        Some(cli::state::load_profile(&state.current_profile))
    } else {
        None
    };

    match args.command {
        Commands::Where => println!(
            "{}",
            force_available(current_profile)
                .get_location()
                .to_string_lossy()
        ),
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
        Commands::Link {
            profile_name,
            overwrite,
        } => {
            let new_profile = if profile_name != "" {
                cli::state::load_profile(&profile_name)
            } else {
                force_available(current_profile)
            };
            new_profile.link(overwrite).expect("Error linking profile!");
            state.current_profile = new_profile.name();
        }
        Commands::Sync { message } => {
            let profile = force_available(current_profile);
            profile
                .sync_with_message(&message)
                .expect("Unable to sync all profiles!");
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
        Commands::Clone {
            url,
            name,
            link,
            overwrite,
        } => {
            let mut work_dir = profiles_dir();
            let mut args = vec!["clone", &url];
            let profile_name = if &name != "" {
                name
            } else {
                url.clone()
                    .split("/")
                    .last()
                    .expect("Unable to parse url!")
                    .to_string()
            };

            args.push(&profile_name);

            Command::new("git")
                .args(args)
                .current_dir(&work_dir)
                .spawn()
                .expect("Error starting git!")
                .wait()
                .expect("Error cloning repository!");

            work_dir.push(&profile_name);

            let profile = Profile::load(&work_dir);
            state.current_profile = profile.name();

            if link {
                profile.link(overwrite).expect("Error linking profile!");
            }
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
        Commands::Clean => {
            let home = cli::state::home_dir();
            let walker = WalkDir::new(&home).into_iter().filter_entry(|entry| {
                // Git repos and node_modules have tons of files and are
                // unlikely to contain dotfiles so this speeds thing up
                // significantly.
                entry.file_name() != ".git" && entry.file_name() != "node_modules"
            });
            let profiles_path = profiles_dir();

            for possible_entry in walker {
                if possible_entry.is_err() {
                    continue;
                }

                let entry = possible_entry.unwrap();
                let path = entry.path();
                if !path.is_symlink() {
                    continue;
                }

                let target = match path.read_link() {
                    Ok(p) => p,
                    Err(_) => continue,
                };

                // If it's not a DFM related symlink ignore it.
                if !target.starts_with(&profiles_path) {
                    continue;
                }

                let printable_path = path.to_string_lossy();
                println!("Checking {}", printable_path);
                let file_exists = target.exists();
                if !file_exists {
                    println!("Link {} is dead removing.", printable_path);
                    fs::remove_file(&path)
                        .expect(format!("Unable to remove file: {}", printable_path).as_ref());
                }
            }
        }
        Commands::Add { files } => {
            let profile = force_available(current_profile);
            let profile_root = profile.get_location();

            for file in files {
                // Get the absolute path of the file so it has $HOME as the
                // prefix.
                let path = Path::new(&file)
                    .canonicalize()
                    .expect(format!("Unable to find file: {}", &file).as_ref());

                let home = cli::state::home_dir()
                    .canonicalize()
                    .expect("Unable to canonicalize home!");

                // Make the path relative to the home directory
                let relative_path = match path.strip_prefix(&home) {
                    Ok(p) => p,
                    Err(_) => {
                        eprintln!("File {} is not in your home directory! If you have a mapping please add it manually.", &file);
                        process::exit(1);
                    }
                };

                // Join the relative directory to the profile root
                let mut target_path = profile_root.clone();
                target_path.push(&relative_path);

                let parent = target_path.parent().unwrap();
                if parent != profile_root {
                    fs::create_dir_all(&parent).expect("Unable to create preceding directories!");
                }

                // Move the file / directory into the profile root
                fs::rename(path, target_path)
                    .expect(format!("Unable to move file: {}", &file).as_ref());
            }

            // Link the profile to create symlinks where files were before.
            profile.link(false).expect("Unable to link profile!");
        }
        Commands::GenCompletions { shell } => match Shell::from_str(&shell, true) {
            Ok(generator) => {
                eprintln!("Generating completion file for {}...", generator);
                let cmd = CLI::command();
                generate(
                    generator,
                    &mut cmd.clone(),
                    cmd.get_name().to_string(),
                    &mut io::stdout(),
                );
            }
            Err(failed) => {
                eprintln!("{} is not a known shell.", failed);
                process::exit(1);
            }
        },
        Commands::External(args) => {
            let plugin_name = format!("dfm-{}", args[0].to_str().unwrap_or_default());
            println!("Calling out to {:?} with {:?}", plugin_name, &args[1..]);
        }
    }

    state.default_save().expect("Unable to save state!");
}
