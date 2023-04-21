use std::{
    env,
    ffi::OsStr,
    fs::{self, File},
    io::{self, BufReader},
    os,
    path::{Path, PathBuf},
    process::{Command, ExitStatus},
    str::FromStr,
};

use serde;
use walkdir::{DirEntry, WalkDir};

use super::hooks::Hooks;

#[derive(Debug, Clone, PartialEq, Eq, serde::Deserialize)]
#[serde(rename_all = "snake_case")]
enum LinkMode {
    Pre,
    Post,
    None,
}

impl Default for LinkMode {
    fn default() -> Self {
        LinkMode::Post
    }
}

fn default_off() -> bool {
    false
}

#[derive(Debug, Clone, serde::Deserialize)]
pub struct DFMConfig {
    #[serde(default)]
    pub location: String,

    #[serde(default = "default_off")]
    prompt_for_commit_message: bool,
    #[serde(default = "default_off")]
    pull_only: bool,
    #[serde(default)]
    link: LinkMode,
    #[serde(default = "Hooks::new")]
    hooks: Hooks,
    #[serde(default = "Vec::new")]
    modules: Vec<DFMConfig>,
}

impl Default for DFMConfig {
    fn default() -> Self {
        DFMConfig {
            prompt_for_commit_message: false,
            pull_only: false,
            link: LinkMode::default(),
            location: "".to_string(),
            hooks: Hooks::new(),
            modules: Vec::new(),
        }
    }
}

impl DFMConfig {
    fn load(file: &Path) -> DFMConfig {
        let fh = File::open(file).expect(
            format!(
                "Unexpected error reading {}",
                file.to_str().unwrap_or(".dfm.yml")
            )
            .as_str(),
        );
        let reader = BufReader::new(fh);
        let mut config: DFMConfig = serde_yaml::from_reader(reader).expect("Malformed .dfm.yml");
        if config.location == "" {
            config.location = file
                .parent()
                .expect("Unexpected error getting profile location!")
                .to_str()
                .expect("Unexpected error turning profile location to a string!")
                .to_string();
        }

        for module in &mut config.modules {
            module.expand();
        }

        config
    }

    fn expand(&mut self) {
        if self.location.starts_with("~") {
            self.location = shellexpand::tilde(&self.location).to_string();
        }
    }
}

#[derive(Debug)]
pub struct Profile {
    pub config: DFMConfig,

    location: PathBuf,
    modules: Vec<Profile>,
}

impl Default for Profile {
    fn default() -> Self {
        Profile {
            config: DFMConfig::default(),
            location: PathBuf::new(),
            modules: Vec::new(),
        }
    }
}

type GitResult = Result<ExitStatus, io::Error>;

fn is_dotfile(entry: &DirEntry) -> bool {
    let filename = entry.file_name().to_str().unwrap_or("");
    // .git files and .dfm.yml are not dotfiles so should be ignored.
    let sys_files = filename == ".dfm.yml" || filename == ".git" || filename == "README.md";
    let is_file = entry.path().is_file();
    !sys_files && is_file
}

// Should return an error
fn remove_if_able(path: &Path, force_remove: bool) -> Option<io::Error> {
    if path.exists() && !path.is_symlink() && !force_remove {
        return Some(io::Error::new(
            io::ErrorKind::AlreadyExists,
            "file exists and is not a symlink, cowardly refusing to remove.",
        ));
    }

    if path.is_dir() {
        return Some(io::Error::new(
            io::ErrorKind::AlreadyExists,
            "directory exists and is not a symlink, cowardly refusing to remove.",
        ));
    }

    if !path.exists() {
        return None;
    }

    println!("Removing {:?}", path);
    fs::remove_file(path).err()
}

impl Profile {
    pub fn load(directory: &Path) -> Profile {
        let path = if directory.starts_with("~") {
            let expanded = shellexpand::tilde(directory.to_str().expect("Invalid directory!"));
            PathBuf::from_str(&expanded).expect("Invalid profile directory!")
        } else {
            directory.to_path_buf().clone()
        };
        let dotdfm = path.join(".dfm.yml");
        if dotdfm.exists() {
            let config = DFMConfig::load(&dotdfm);
            return Profile::from_config(config);
        }

        Profile::default()
    }

    pub fn from_config(config: DFMConfig) -> Profile {
        let modules = config
            .modules
            .iter()
            .map(|cfg| Profile::from_config_ref(cfg))
            .collect();
        let location = PathBuf::from_str(&config.location)
            .expect("Unable to convert config location into a path!");

        Profile {
            config: config,
            location: location,
            modules: modules,
        }
    }

    fn from_config_ref(config: &DFMConfig) -> Profile {
        Profile::from_config(config.clone())
    }

    pub fn name(&self) -> String {
        match self.location.file_name() {
            None => "".to_string(),
            Some(basename) => basename.to_string_lossy().to_string(),
        }
    }

    // TODO: hooks
    pub fn link(&self) -> Result<(), io::Error> {
        for profile in self
            .modules
            .iter()
            .filter(|p| p.config.link == LinkMode::Pre)
        {
            profile.link()?;
        }

        let walker = WalkDir::new(&self.location)
            .min_depth(1)
            .into_iter()
            .filter_entry(is_dotfile);

        let home = PathBuf::from(env::var("HOME").unwrap_or("".to_string()));
        for possible_entry in walker {
            let entry = match possible_entry {
                Ok(d) => d,
                Err(_) => continue,
            };
            let file = entry.path();
            let relative_path = file.strip_prefix(&self.location).unwrap();
            let target_path = home.join(relative_path);
            println!(
                "Link {} -> {}",
                target_path.to_string_lossy(),
                file.to_string_lossy()
            );

            if let Some(err) = remove_if_able(&target_path, false) {
                if err.kind() == io::ErrorKind::AlreadyExists {
                    continue;
                }

                return Err(err);
            }

            os::unix::fs::symlink(file, target_path)?;
        }

        for profile in self
            .modules
            .iter()
            .filter(|p| p.config.link == LinkMode::Post)
        {
            profile.link()?;
        }

        Ok(())
    }

    fn git<I, S>(&self, args: I) -> GitResult
    where
        I: IntoIterator<Item = S>,
        S: AsRef<OsStr>,
    {
        Command::new("git")
            .args(args)
            .current_dir(&self.location)
            .spawn()?
            .wait()
    }

    pub fn status(&self) -> GitResult {
        self.git(["status"])
    }
}
