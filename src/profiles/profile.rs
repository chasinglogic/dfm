use std::{
    env,
    ffi::OsStr,
    fs::{self, File},
    io::{self, Write},
    os,
    path::{Path, PathBuf},
    process::{Command, ExitStatus},
    str::FromStr,
};

use super::config::DFMConfig;

use walkdir::{DirEntry, WalkDir};

#[derive(Debug)]
pub struct Profile {
    config: DFMConfig,

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
    let is_sys_file = filename == ".dfm.yml" || filename == ".git" || filename == "README.md";
    let is_file = entry.path().is_file();
    !is_sys_file && is_file
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

        let mut profile = Profile::default();
        profile.config.location = path.to_string_lossy().to_string();
        profile.location = path;
        profile
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

    pub fn link(&self) -> Result<(), io::Error> {
        for profile in self
            .modules
            .iter()
            .filter(|p| p.config.link == LinkMode::Pre)
        {
            profile.link()?;
        }

        self.run_hook("before_link")?;

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

        self.run_hook("after_link")?;

        for profile in self
            .modules
            .iter()
            .filter(|p| p.config.link == LinkMode::Post)
        {
            profile.link()?;
        }

        Ok(())
    }

    pub fn git<I, S>(&self, args: I) -> GitResult
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

    pub fn init(&self) -> Result<(), io::Error> {
        self.git(["init"])?;

        let mut dotdfm = self.location.clone();
        dotdfm.push(".dfm.yml");
        let fh = &mut File::create(&dotdfm)?;
        // TODO: Embed a hardcoded default config with documentation comments
        // and good formatting in the binary and use it here.
        let content = serde_yaml::to_string(&self.config)
            .map_err(|e| io::Error::new(io::ErrorKind::InvalidData, e.to_string()))?;
        fh.write_all(content.as_bytes())?;

        self.git(["add", ".dfm.yml"])?;
        self.git(["commit", "-m", "initial commit"])?;

        Ok(())
    }

    pub fn run_hook(&self, hook_name: &str) -> Result<(), io::Error> {
        self.config.hooks.run_hook(hook_name, &self.location)
    }

    pub fn get_location(&self) -> PathBuf {
        self.location.clone()
    }
}
