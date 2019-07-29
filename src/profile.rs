use serde::Deserialize;
use std::env;
use std::fmt;
use std::fs::File;
use std::io;
use std::io::prelude::*;
use std::path::{Path, PathBuf};

use log::debug;

use crate::hooks::{HookConfig, Hooks};
use crate::mapping::{MappingConfig, Mappings};
use crate::repo::Repo;

pub struct Error {
    profile: Profile,
    error: io::Error,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}: {}", self.profile.repo.path.display(), self.error)
    }
}

fn default_off() -> bool {
    false
}

#[derive(Deserialize)]
pub struct ProfileConfig {
    #[serde(default = "String::new")]
    target_dir: String,
    #[serde(default = "String::new")]
    location: String,
    #[serde(default = "default_off")]
    pull_only: bool,
    #[serde(default = "String::new")]
    link: String,

    #[serde(default = "String::new")]
    commit_msg: String,
    #[serde(default = "default_off")]
    prompt_for_commit_message: bool,
    hooks: Option<HookConfig>,
    mappings: Option<Vec<MappingConfig>>,
    modules: Option<Vec<ProfileConfig>>,
}

impl ProfileConfig {
    pub fn default() -> ProfileConfig {
        ProfileConfig {
            location: String::new(),
            link: "post".to_string(),
            pull_only: false,

            target_dir: env::var("HOME").unwrap_or_default(),
            commit_msg: env::var("DFM_COMMIT_MSG").unwrap_or_default(),
            prompt_for_commit_message: false,
            hooks: None,
            mappings: None,
            modules: None,
        }
    }
}

#[derive(Debug, Clone)]
pub struct Profile {
    pub hooks: Hooks,
    pub repo: Repo,
    pub target_dir: PathBuf,
    pub pull_only: bool,

    commit_msg: String,
    link_when: String,
    mappings: Mappings,
    modules: Vec<Profile>,
    prompt_for_commit_message: bool,
}

impl Profile {
    fn from(profile_dir: &Path, config: ProfileConfig) -> Profile {
        let target_dir = if config.target_dir != "" {
            PathBuf::from(config.target_dir)
        } else {
            PathBuf::from(env::var("HOME").unwrap_or_default())
        };
        let hooks = Hooks::from(config.hooks.unwrap_or_default());
        let mappings = Mappings::from(config.mappings.unwrap_or_default());
        let modules = config
            .modules
            .unwrap_or_default()
            .drain(..)
            .map(|cfg| {
                let location = cfg
                    .location
                    .clone()
                    .replace("~", &env::var("HOME").unwrap_or_default());
                let path = Path::new(&location);
                Profile::from(path, cfg)
            })
            .collect();
        Profile {
            mappings,
            modules,
            target_dir,
            hooks,

            commit_msg: config.commit_msg,
            link_when: config.link,
            prompt_for_commit_message: config.prompt_for_commit_message,
            pull_only: config.pull_only,
            repo: Repo::new(profile_dir),
        }
    }

    pub fn set_commit_msg(&mut self, msg: &str) {
        self.commit_msg = msg.to_string();
    }

    pub fn set_prompt(&mut self, prompt: bool) {
        self.prompt_for_commit_message = prompt;
    }

    pub fn load(profile_dir: &Path) -> Result<Profile, io::Error> {
        let dir = profile_dir.to_path_buf();
        let mut cfg_file = dir.clone();
        cfg_file.push(".dfm.yml");
        let config = if cfg_file.exists() && cfg_file.is_file() {
            let mut cfg = File::open(&cfg_file)?;
            let mut contents = String::new();
            cfg.read_to_string(&mut contents)?;
            match serde_yaml::from_str(&contents) {
                Ok(c) => c,
                Err(e) => return Err(io::Error::new(io::ErrorKind::Other, format!("{}", e))),
            }
        } else {
            ProfileConfig::default()
        };

        Ok(Profile::from(profile_dir, config))
    }

    pub fn sync(&self) -> Result<(), Error> {
        debug!("running pre sync hooks");
        self.run_hooks("pre", "sync")?;

        for module in self.modules.iter().filter(|p| p.link_when == "pre") {
            debug!("syncing module: {:?}", module);
            module.sync()?;
        }

        println!("\n{}:", self.repo.path.display());
        let input: String;
        // TODO: show a diff when prompting
        let msg: &str = if self.repo.is_dirty() && self.prompt_for_commit_message {
            self.git(&["diff"])?;
            print!("Commit msg: ");
            io::stdout().flush().expect("unable to flush stdout");
            input = read!("{}\n");
            &input
        } else {
            &self.commit_msg
        };

        if let Err(e) = self.repo.sync(msg, self.pull_only) {
            return Err(Error {
                profile: self.clone(),
                error: e,
            });
        };

        for module in self.modules.iter().filter(|p| p.link_when != "pre") {
            debug!("syncing module: {:?}", module);
            module.sync()?;
        }

        debug!("running post sync hooks");
        self.run_hooks("post", "sync")?;
        Ok(())
    }

    pub fn init(&self) -> Result<(), Error> {
        self.git(&["init"])
    }

    pub fn link(&self, overwrite: bool) -> Result<(), Error> {
        if self.link_when == "none" {
            return Ok(());
        }

        self.run_hooks("pre", "link")?;
        let target_dir = Path::new(&self.target_dir);
        for module in self.modules.iter().filter(|p| p.link_when == "pre") {
            debug!("linking module: {}", module.repo.path.display());
            module.link(overwrite)?;
        }

        debug!("target directory: {}", target_dir.display());
        // doesn't link modules?
        let links = match self.mappings.link(&self.repo.path, &target_dir, overwrite) {
            Ok(links) => links,
            Err(e) => {
                return Err(Error {
                    profile: self.clone(),
                    error: e,
                })
            }
        };
        for link in links.iter() {
            debug!("link {} => {}", link.src.display(), link.dst.display());
            if let Err(e) = link.link() {
                return Err(Error {
                    profile: self.clone(),
                    error: e,
                });
            };
        }

        for module in self.modules.iter().filter(|p| p.link_when != "pre") {
            debug!("linking module: {}", module.repo.path.display());
            module.link(overwrite)?;
        }

        self.run_hooks("post", "link")?;
        Ok(())
    }

    fn run_hooks(&self, when: &str, name: &str) -> Result<(), Error> {
        let res = if when == "pre" {
            self.hooks.pre(&self.repo.path, name)
        } else {
            self.hooks.post(&self.repo.path, name)
        };

        match res {
            Err(e) => Err(Error {
                profile: self.clone(),
                error: e,
            }),
            Ok(_) => Ok(()),
        }
    }

    fn git(&self, cmd: &[&str]) -> Result<(), Error> {
        if let Err(e) = self.repo.git(cmd) {
            return Err(Error {
                profile: self.clone(),
                error: e,
            });
        }

        Ok(())
    }
}
