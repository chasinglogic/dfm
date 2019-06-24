use serde::Deserialize;
use std::env;
use std::fs::File;
use std::io;
use std::io::prelude::*;
use std::path::{Path, PathBuf};

use crate::hooks::{HookConfig, Hooks};
use crate::mapping::{MappingConfig, Mappings};
use crate::repo::Repo;

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

            target_dir: env::var("HOME").unwrap_or("".to_string()),
            commit_msg: env::var("DFM_COMMIT_MSG").unwrap_or("".to_string()),
            prompt_for_commit_message: false,
            hooks: None,
            mappings: None,
            modules: None,
        }
    }
}

pub struct Profile {
    pub hooks: Hooks,
    pub repo: Repo,

    commit_msg: String,
    link_when: String,
    mappings: Mappings,
    modules: Vec<Profile>,
    prompt_for_commit_message: bool,
    pull_only: bool,
    target_dir: PathBuf,
}

impl Profile {
    fn from(profile_dir: &Path, mut config: ProfileConfig) -> Profile {
        let target_dir = PathBuf::from(config.target_dir);
        let hooks = Hooks::from(config.hooks.unwrap_or(HookConfig::new()));
        let mappings = Mappings::from(config.mappings.unwrap_or(Vec::new()));
        let modules = config
            .modules
            .unwrap_or(Vec::new())
            .drain(..)
            .map(|cfg| {
                let location = cfg.location.clone();
                Profile::from(Path::new(&location), cfg)
            })
            .collect();
        Profile {
            commit_msg: config.commit_msg,
            hooks: hooks,
            link_when: config.link,
            mappings: mappings,
            modules: modules,
            prompt_for_commit_message: config.prompt_for_commit_message,
            pull_only: config.pull_only,
            target_dir: target_dir,
            repo: Repo::new(profile_dir),
        }
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

    pub fn sync(&self) -> Result<(), io::Error> {
        self.hooks.pre(&self.repo.path, "sync")?;

        for module in self.modules.iter().filter(|p| p.link_when == "pre") {
            module.sync()?;
        }

        let input: String;
        // TODO: show a diff when prompting
        let msg: &str = if self.repo.is_dirty() && self.prompt_for_commit_message {
            print!("Commit msg: ");
            io::stdout().flush();
            input = read!("{}\n");
            &input
        } else {
            &self.commit_msg
        };
        self.repo.sync(msg, self.pull_only)?;

        for module in self.modules.iter().filter(|p| p.link_when != "pre") {
            module.sync()?;
        }

        self.hooks.post(&self.repo.path, "sync")?;
        Ok(())
    }

    pub fn link(&self) -> Result<(), io::Error> {
        self.hooks.pre(&self.repo.path, "link")?;
        let target_dir = Path::new(&self.target_dir);
        let links = self.mappings.link(&self.repo.path, target_dir)?;
        for link in links.iter() {
            link.link()?;
        }
        self.hooks.post(&self.repo.path, "link")?;
        Ok(())
    }
}
