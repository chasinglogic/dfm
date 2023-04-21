use std::{
    fs::File,
    io::BufReader,
    path::{Path, PathBuf},
    str::FromStr,
};

use serde;

use super::hooks::Hooks;

#[derive(Debug, Clone, serde::Deserialize)]
#[serde(rename_all = "snake_case")]
enum LinkMode {
    Default,
    None,
    Pre,
    Post,
}

impl Default for LinkMode {
    fn default() -> Self {
        LinkMode::Default
    }
}

fn default_off() -> bool {
    false
}

#[derive(Debug, Clone, serde::Deserialize)]
pub struct DFMConfig {
    #[serde(default = "default_off")]
    prompt_for_commit_message: bool,
    #[serde(default = "default_off")]
    pull_only: bool,
    #[serde(default)]
    link: LinkMode,
    #[serde(default)]
    pub location: String,
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
}
