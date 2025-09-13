use crate::utils;
use std::{fs::File, io::BufReader, path::Path, path::PathBuf};

use super::hooks::Hooks;
use super::mapping::Mapping;

#[derive(Default, Debug, Clone, PartialEq, Eq, serde::Deserialize, serde::Serialize)]
#[serde(rename_all = "snake_case")]
pub enum LinkMode {
    Pre,
    #[default]
    Post,
    None,
}

fn default_off() -> bool {
    false
}

#[derive(Debug, Clone, serde::Deserialize, serde::Serialize)]
pub struct DFMConfig {
    #[serde(default, alias = "repository")]
    pub repo: String,
    #[serde(default)]
    pub location: String,
    #[serde(default = "Hooks::new")]
    pub hooks: Hooks,

    #[serde(default)]
    pub root_dir: Option<String>,

    #[serde(default)]
    pub clone_flags: Vec<String>,

    #[serde(default = "default_off")]
    pub prompt_for_commit_message: bool,
    #[serde(default = "default_off")]
    pub pull_only: bool,
    #[serde(default)]
    pub link: LinkMode,
    #[serde(default = "Vec::new")]
    pub modules: Vec<DFMConfig>,
    pub mappings: Option<Vec<Mapping>>,
}

impl Default for DFMConfig {
    fn default() -> Self {
        DFMConfig {
            prompt_for_commit_message: false,
            pull_only: false,
            link: LinkMode::default(),
            repo: "".to_string(),
            location: "".to_string(),
            root_dir: None,
            hooks: Hooks::new(),
            modules: Vec::new(),
            mappings: None,
            clone_flags: Vec::new(),
        }
    }
}

impl DFMConfig {
    pub fn load(file: &Path) -> DFMConfig {
        let fh = File::open(file).unwrap_or_else(|_| {
            panic!(
                "Unexpected error reading {}",
                file.to_str().unwrap_or(".dfm.yml")
            )
        });
        let reader = BufReader::new(fh);
        let mut config: DFMConfig = serde_yaml::from_reader(reader).expect("Malformed .dfm.yml");
        if config.location.is_empty() {
            config.location = file
                .parent()
                .expect("Unexpected error getting profile location!")
                .to_str()
                .expect("Unexpected error turning profile location to a string!")
                .to_string();
        }

        for module in &mut config.modules {
            module.expand_module();
        }

        config
    }

    pub fn get_location(&self) -> PathBuf {
        let mut location = utils::expand_path(&self.location);

        if let Some(ref new_root) = self.root_dir {
            location.push(new_root);
        }

        location
    }

    fn expand_module(&mut self) {
        if self.location.is_empty() {
            let name = self
                .repo
                .split('/')
                .next_back()
                .expect("A module must define a repository!")
                .replace(".git", "");

            let mut module_dir = crate::cli::state::modules_dir();
            module_dir.push(name);

            self.location = module_dir.to_string_lossy().to_string();
        } else {
            self.location = shellexpand::tilde(&self.location).to_string();
        }
    }
}
