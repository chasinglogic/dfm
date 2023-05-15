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

use super::config::{DFMConfig, LinkMode};

use text_io::read;
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
    !is_sys_file
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
        let modules: Vec<Profile> = config
            .modules
            .iter()
            .map(|cfg| Profile::from_config_ref(cfg))
            .collect();

        for module in modules.iter() {
            if !module.get_location().exists() {
                module.download();
            }
        }

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

    fn download(&self) {
        Command::new("git")
            .args([
                "clone",
                &self.config.repo,
                self.get_location().to_str().expect("Unexpected error!"),
            ])
            .spawn()
            .expect("Unable to start git clone!")
            .wait()
            .expect(format!("Unable to clone module! {}", self.config.repo).as_str());
    }

    pub fn name(&self) -> String {
        match self.location.file_name() {
            None => "".to_string(),
            Some(basename) => basename.to_string_lossy().to_string(),
        }
    }

    pub fn is_dirty(&self) -> bool {
        let mut proc = Command::new("git");
        proc.args(["status", "--porcelain"]);
        proc.current_dir(&self.location);

        match proc.output() {
            Ok(output) => output.stdout != "".as_bytes(),
            Err(_) => false,
        }
    }

    pub fn has_origin(&self) -> bool {
        let mut proc = Command::new("git");
        proc.args(["remote", "-v"]);
        proc.current_dir(&self.location);

        match proc.output() {
            Ok(output) => {
                let remotes = String::from_utf8(output.stdout).unwrap_or("".to_string());
                remotes.contains("origin")
            }
            Err(_) => false,
        }
    }

    pub fn branch_name(&self) -> String {
        let mut proc = Command::new("git");
        proc.args(["rev-parse", "--abbrev-ref", "HEAD"]);
        proc.current_dir(&self.location);

        match proc.output() {
            Ok(output) => {
                let branch = String::from_utf8(output.stdout).unwrap_or("".to_string());
                branch.trim().to_string()
            }
            Err(_) => "main".to_string(),
        }
    }

    pub fn sync(&self) -> Result<(), io::Error> {
        self.sync_with_message("")
    }

    pub fn sync_with_message(&self, commit_msg: &str) -> Result<(), io::Error> {
        let is_dirty = self.is_dirty();
        let has_origin = self.has_origin();
        let branch_name = self.branch_name();

        if is_dirty {
            let msg = if self.config.prompt_for_commit_message && commit_msg.is_empty() {
                self.git(["--no-pager", "diff"])?;
                print!("Commit message: ");
                read!("{}\n")
            } else if !commit_msg.is_empty() {
                commit_msg.to_string()
            } else {
                "Dotfiles managed by DFM! https://github.com/chasinglogic/dfm".to_string()
            };

            self.run_hook("before_sync")?;
            self.git(["add", "--all"])?;
            self.git(["commit", "-m", &msg])?;
        }

        if has_origin {
            self.git(["pull", "--rebase", "origin", &branch_name])?;
        }

        if is_dirty && has_origin {
            self.git(["push", "origin", &branch_name])?;
            self.run_hook("after_sync")?;
        }

        for profile in &self.modules {
            profile.sync()?;
        }

        Ok(())
    }

    pub fn link(&self, overwrite_existing_files: bool) -> Result<(), io::Error> {
        for profile in self
            .modules
            .iter()
            .filter(|p| p.config.link == LinkMode::Pre)
        {
            profile.link(overwrite_existing_files)?;
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
            if !file.is_file() {
                continue;
            }

            let relative_path = file.strip_prefix(&self.location).unwrap();
            let target_path = home.join(relative_path);
            println!(
                "Link {} -> {}",
                target_path.to_string_lossy(),
                file.to_string_lossy()
            );

            if let Some(err) = remove_if_able(&target_path, overwrite_existing_files) {
                if err.kind() == io::ErrorKind::AlreadyExists {
                    eprintln!("{}", err);
                    continue;
                }

                return Err(err);
            }

            if let Some(path) = target_path.parent() {
                if !path.exists() {
                    fs::create_dir_all(path)?;
                }
            }

            os::unix::fs::symlink(file, target_path)?;
        }

        self.run_hook("after_link")?;

        for profile in self
            .modules
            .iter()
            .filter(|p| p.config.link == LinkMode::Post)
        {
            profile.link(overwrite_existing_files)?;
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
