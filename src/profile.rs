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
use crate::util;

pub struct Error {
    profile_path: PathBuf,
    error: io::Error,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}: {}", self.profile_path.display(), self.error)
    }
}

fn default_off() -> bool {
    false
}

#[derive(Deserialize)]
pub struct ProfileConfig {
    #[serde(default = "String::new")]
    repo: String,
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
            repo: String::new(),
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
    fn update_config(mut self, cfg: ProfileConfig) -> Profile {
        if cfg.target_dir != "" {
            self.target_dir = PathBuf::from(cfg.target_dir);
        }

        self.pull_only = cfg.pull_only;
        self.commit_msg = cfg.commit_msg;
        self.prompt_for_commit_message = cfg.prompt_for_commit_message;
        self.link_when = cfg.link;
        self
    }

    fn from(profile_dir: &Path, config: ProfileConfig) -> Result<Profile, Error> {
        let target_dir = if config.target_dir != "" {
            PathBuf::from(config.target_dir)
        } else {
            PathBuf::from(env::var("HOME").unwrap_or_default())
        };
        let hooks = Hooks::from(config.hooks.unwrap_or_default());
        let mappings = Mappings::from(config.mappings.unwrap_or_default());
        let module_configs = config.modules.unwrap_or_default();

        let mut modules = Vec::with_capacity(module_configs.len());
        for cfg in module_configs {
            let location = cfg
                .location
                .clone()
                .replace("~", &env::var("HOME").unwrap_or_default());

            let path = if location == "" {
                let name = util::default_profile_name(&cfg.repo);
                let mut sd = util::profile_storage_dir(&util::cfg_dir(None));
                sd.push(name);
                sd
            } else {
                PathBuf::from(&location)
            };

            let profile = Profile::load(&path)?.update_config(cfg);
            modules.push(profile);
        }

        Ok(Profile {
            mappings,
            modules,
            target_dir,
            hooks,

            commit_msg: config.commit_msg,
            link_when: config.link,
            prompt_for_commit_message: config.prompt_for_commit_message,
            pull_only: config.pull_only,
            repo: match Repo::new(profile_dir, &config.repo) {
                Ok(r) => r,
                Err(e) => {
                    return Err(Error {
                        profile_path: profile_dir.to_path_buf(),
                        error: e,
                    })
                }
            },
        })
    }

    pub fn set_commit_msg(&mut self, msg: &str) {
        self.commit_msg = msg.to_string();
    }

    pub fn set_prompt(&mut self, prompt: bool) {
        self.prompt_for_commit_message = prompt;
    }

    pub fn load(profile_dir: &Path) -> Result<Profile, Error> {
        let dir = profile_dir.to_path_buf();
        let mut cfg_file = dir.clone();
        cfg_file.push(".dfm.yml");
        let config = if cfg_file.exists() && cfg_file.is_file() {
            let mut cfg = match File::open(&cfg_file) {
                Ok(f) => f,
                Err(e) => {
                    return Err(Error {
                        profile_path: profile_dir.to_path_buf(),
                        error: e,
                    })
                }
            };
            let mut contents = String::new();
            if let Err(e) = cfg.read_to_string(&mut contents) {
                return Err(Error {
                    profile_path: profile_dir.to_path_buf(),
                    error: e,
                });
            }
            match serde_yaml::from_str(&contents) {
                Ok(c) => c,
                Err(e) => {
                    return Err(Error {
                        profile_path: profile_dir.to_path_buf(),
                        error: io::Error::new(
                            io::ErrorKind::Other,
                            format!("Invalid YAML in .dfm.yml: {}", e),
                        ),
                    })
                }
            }
        } else {
            ProfileConfig::default()
        };

        Ok(Profile::from(profile_dir, config)?)
    }

    pub fn sync(&self) -> Result<(), Error> {
        debug!("running pre sync hooks");
        self.run_hooks("pre", "sync")?;

        for module in self.modules.iter().filter(|p| p.link_when == "pre") {
            debug!("syncing module: {:?}", module);
            module.sync()?;
        }

        println!("{}:", self.repo.path.display());
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
                profile_path: self.repo.path.clone(),
                error: e,
            });
        };

        // Print a newline after sync to separate git output from repo path.
        println!();

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

        debug!(
            "running pre link hooks for profile: {}",
            self.repo.path.display()
        );
        self.run_hooks("pre", "link")?;
        let target_dir = Path::new(&self.target_dir);
        for module in self.modules.iter().filter(|p| p.link_when == "pre") {
            debug!("linking module: {}", module.repo.path.display());
            module.link(overwrite)?;
        }

        debug!("target directory: {}", target_dir.display());
        let links = match self.mappings.link(&self.repo.path, &target_dir, overwrite) {
            Ok(links) => links,
            Err(e) => {
                return Err(Error {
                    profile_path: self.repo.path.clone(),
                    error: e,
                })
            }
        };
        for link in links.iter() {
            debug!("link {} => {}", link.src.display(), link.dst.display());
            if let Err(e) = link.link() {
                return Err(Error {
                    profile_path: self.repo.path.clone(),
                    error: e,
                });
            };
        }

        for module in self.modules.iter().filter(|p| p.link_when != "pre") {
            debug!("linking module: {}", module.repo.path.display());
            module.link(overwrite)?;
        }

        debug!(
            "running post link hooks for profile: {}",
            self.repo.path.display()
        );
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
                profile_path: self.repo.path.clone(),
                error: e,
            }),
            Ok(_) => Ok(()),
        }
    }

    fn git(&self, cmd: &[&str]) -> Result<(), Error> {
        if let Err(e) = self.repo.git(cmd) {
            return Err(Error {
                profile_path: self.repo.path.clone(),
                error: e,
            });
        }

        Ok(())
    }
}
